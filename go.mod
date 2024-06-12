module github.com/cern-eos/eosxd-csi

go 1.22.0

toolchain go1.22.2

require (
	github.com/container-storage-interface/spec v1.9.0
	github.com/kubernetes-csi/csi-lib-utils v0.18.1
	github.com/moby/sys/mountinfo v0.7.1
	google.golang.org/grpc v1.64.0
	k8s.io/apimachinery v0.30.1
	k8s.io/klog/v2 v2.120.1
	k8s.io/mount-utils v0.30.1
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	k8s.io/utils v0.0.0-20240502163921-fe8a2dddb1d0 // indirect
)
