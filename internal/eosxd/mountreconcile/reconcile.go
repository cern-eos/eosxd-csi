// Copyright CERN.
//
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package mountreconcile

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/cern-eos/eosxd-csi/internal/log"
	"github.com/cern-eos/eosxd-csi/internal/mountutils"

	"github.com/moby/sys/mountinfo"
)

type Opts struct {
	Frequency time.Duration
}

func RunBlocking(o *Opts) error {
	t := time.NewTicker(o.Frequency)

	doReconcile := func() {
		log.Tracef("Reconciling /cvmfs")
		if err := reconcile(); err != nil {
			log.Errorf("Failed to reconcile /cvmfs: %v", err)
		}
	}

	// Run at start so that broken mounts after nodeplugin Pod
	// restart are cleaned up.
	doReconcile()

	for {
		select {
		case <-t.C:
			doReconcile()
		}
	}
}

// List EOSxd mounts that the kernel knows about.
// Maps EOS instance to slice of mountpoints.
func getKnownEosMountsGroupedByInstance() (map[string][]string, error) {
	eosMountInfos, err := mountinfo.GetMounts(func(info *mountinfo.Info) (skip, stop bool) {
		skip = info.FSType != "fuse" || !strings.HasPrefix(info.Mountpoint, "/eos/")
		return
	})
	if err != nil {
		return nil, err
	}

	eosMountsByInstances := make(map[string][]string)
	for _, info := range eosMountInfos {
		eosMountsByInstances[info.Source] = append(eosMountsByInstances[info.Source], info.Mountpoint)
	}

	return eosMountsByInstances, nil
}

func isAllNum(s string) bool {
	sz := len(s)
	for i := 0; i < sz; i++ {
		c := s[i]
		if c < '0' && c > '9' {
			return false
		}
	}

	return true
}

// List EOSxd mounts that are served by eosxd processes.
// Maps EOS instance to slice of mountpoints.
func getRealEosMountsGroupedByInstance() (map[string][]string, error) {
	procDir, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("failed to read /proc: %v", err)
	}

	// Filter procs that are running eosxd.

	var eosCmdlines [][]byte

	for _, proc := range procDir {
		if !proc.IsDir() || !isAllNum(proc.Name()) {
			continue
		}

		cmdlinePath := path.Join("/proc", proc.Name(), "cmdline")
		cmdline, err := os.ReadFile(cmdlinePath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			log.Errorf("Failed to read cmdline from %s: %v", cmdlinePath, err)

			return nil, err
		}

		progNameEndIdx := bytes.IndexByte(cmdline, 0)
		if progNameEndIdx == -1 {
			continue
		}

		if !bytes.HasSuffix(cmdline[:progNameEndIdx], []byte("/eosxd")) {
			continue
		}

		eosCmdlines = append(eosCmdlines, cmdline)
	}

	// Retrieve mountpoints and EOS instance names from the cmdlines.

	eosMountsByInstances := make(map[string][]string)

	for _, cmdline := range eosCmdlines {
		// eosxd command line is expected to be in consisted
		// of following parts delimited by null characters:
		//
		//   * <eosxd path> (usually /usr/bin/eosxd)
		//   * <mountpoint>
		//   * "-o" (mount options flag)
		//   * <mount flags> (usually "rw,fsname=<EOS instance>")
		//   * "-oautofs" (extra flag added by autofs)
		//   * (empty)

		parts := bytes.Split(cmdline, []byte{0})

		if len(parts) != 6 {
			log.Errorf("Failed to parse cmdline %s: unexpected number of arguments: got %d, wanted 5",
				cmdline, len(parts),
			)
			continue
		}

		const fsnameFlag = "fsname="

		var (
			instanceName string
			mountpoint   = string(parts[1])
		)

		// Extract fsname=<EOS instance> from the mount flags.

		flags := bytes.Split(parts[3], []byte{','})
		for i := range flags {
			if bytes.HasPrefix(flags[i], []byte(fsnameFlag)) {
				instanceName = string(flags[i][len(fsnameFlag):])
				break
			}
		}

		if instanceName == "" {
			log.Errorf("Failed to extract fsname from cmdline %s", cmdline)
			continue
		}

		// Store the results.

		eosMountsByInstances[string(instanceName)] = append(eosMountsByInstances[string(instanceName)], mountpoint)
	}

	return eosMountsByInstances, nil
}

func reconcile() error {
	// List EOS mount info from two sources:
	// * mounts that the kernel knows about (/proc/self/mountinfo),
	// * listing /proc and searching for eosxd processes.

	knownEosMounts, err := getKnownEosMountsGroupedByInstance()
	if err != nil {
		return err
	}

	realEosMounts, err := getRealEosMountsGroupedByInstance()
	if err != nil {
		return err
	}

	log.Tracef("EOSxd mounts known to kernel: %v", knownEosMounts)
	log.Tracef("eosxd processes running for mounts: %v", realEosMounts)

	// Make sure there is an eosxd process running for each known mount.
	// If not, reconcile by unmounting. autofs will automatically remount when needed.

	for instance, mountpoints := range knownEosMounts {
		if _, ok := realEosMounts[instance]; ok {
			// eosxd process exists, move to the next instance/mounts.
			continue
		}

		log.Infof("Mountpoints %v for instance %s are corrupted, unmounting", mountpoints, instance)

		for i := range mountpoints {
			if err := mountutils.Unmount(mountpoints[i]); err != nil {
				log.Errorf("Failed to unmount %s during mount reconciliation", mountpoints[i])
			}
		}
	}

	return nil
}
