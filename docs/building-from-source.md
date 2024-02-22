# Building EOSxd CSI from source

There are pre-built container images available from [eosxd-csi repository in CERN registry](https://registry.cern.ch/harbor/projects/2335/repositories/eosxd-csi/artifacts-tab) (`docker pull registry.cern.ch/kubernetes/eosxd-csi`). If however you need to build EOSxd CSI from source, you can follow this guide.

EOSxd CSI is written in Go. Make sure you have Go compiler and other related build tools installed before continuing.

Clone [github.com/cern-eos/eosxd-csi](https://github.com/cern-eos/eosxd-csi) repository:
```bash
git clone https://github.com/cern-eos/eosxd-csi.git
cd eosxd-csi
```

There are different build targets available in the provided Makefile. To build only the EOSxd CSI executables, run following command:
```bash
make
```

After building successfully, the resulting executable files can be found in `bin/<Platform>-<Architecture>` (e.g. `bin/linux-amd64`).

You can also build container images. By default, Docker is used for building. To build a container image using e.g. Podman, run following command:
```bash
TARGETS=linux/amd64 IMAGE_BUILD_TOOL=podman make image
```
