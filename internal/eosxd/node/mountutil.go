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

package node

import (
	goexec "os/exec"

	"github.com/cern-eos/eosxd-csi/internal/exec"
	"github.com/cern-eos/eosxd-csi/internal/mountutils"
)

func bindMount(from, to string) error {
	_, err := exec.CombinedOutput(goexec.Command("mount", "--bind", from, to))
	return err
}

func slaveRecursiveBind(from, to string) error {
	_, err := exec.CombinedOutput(goexec.Command(
		"mount",
		from,
		to,

		// We bindmount recursively in order to retain any
		// existing EOSxd mounts inside of the autofs root.
		"--rbind",

		// We expect the autofs root in /eos to be already marked
		// as shared, making it possible to send and receive mount
		// and unmount events between bindmounts. We need to make event
		// propagation one-way only (from autofs root to bindmounts)
		// however, because, when unmounting, we do so recursively, and
		// this would then mean attempting to unmount autofs-EOS mounts
		// in the rest of the bindmounts (used by other Pods on the node
		// that also use EOS), which is not desirable of course.
		"--make-slave",
	))

	return err
}

func makeRecursivePrivateMount(mountpoint string) error {
	_, err := exec.CombinedOutput(goexec.Command("mount", mountpoint, "--make-rprivate"))
	return err
}

func recursiveUnmount(mountpoint string) error {
	// We need recursive unmount because there are live mounts inside the bindmount.
	// Unmounting only the upper autofs mount would result in EBUSY.
	return mountutils.Unmount(mountpoint, "--recursive")
}
