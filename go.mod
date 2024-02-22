module github.com/cern-eos/eosxd-csi

go 1.21

toolchain go1.21.5

require (
	github.com/container-storage-interface/spec v1.9.0
	github.com/kubernetes-csi/csi-lib-utils v0.17.0
	github.com/moby/sys/mountinfo v0.7.1
	google.golang.org/grpc v1.62.0
	k8s.io/apimachinery v0.29.1
	k8s.io/klog/v2 v2.120.1
	k8s.io/mount-utils v0.29.1
)

require (
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240123012728-ef4313101c80 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
	k8s.io/utils v0.0.0-20230726121419-3b25d923346b // indirect
)
