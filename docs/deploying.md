# Deploying EOSxd CSI driver in Kubernetes

EOSxd CSI deployment consists of a DaemonSet node plugin that handles node local mount-unmount operations, and ConfigMaps storing EOSxd client configuration.

Cluster administrators may deploy EOSxd CSI manually using the provided Kubernetes manifests, or by installing eosxd-csi Helm chart.

After successful deployment, you can try examples in [../example/](../example/).

## Deployment with Helm chart

Helm chart can be installed from CERN registry:

```bash
helm install eosxd-csi oci://registry.cern.ch/kubernetes/charts/eosxd-csi --version <Chart tag>
```

Some chart values may need to be customized to suite your EOSxd environment. Please consult the documentation in [../deployments/helm/README.md](../deployments/helm/README.md) to see available values.

## Verifying the deployment

After successful deployment, you should see similar output from `kubectl get all -l app=eosxd-csi`:

```
$ kubectl -n eos get all -l app=eosxd-csi
NAME                                              READY   STATUS    RESTARTS   AGE
pod/eosxd-csi-controllerplugin-7dbdfd6569-79n77   2/2     Running   0          22s
pod/eosxd-csi-nodeplugin-94sbp                    4/4     Running   0          22s

NAME                                  DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR   AGE
daemonset.apps/eosxd-csi-nodeplugin   1         1         1       1            1           <none>          22s

NAME                                         READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/eosxd-csi-controllerplugin   1/1     1            1           23s

NAME                                                    DESIRED   CURRENT   READY   AGE
replicaset.apps/eosxd-csi-controllerplugin-7dbdfd6569   1         1         1       23s
```

## csi-driver command line arguments

EOSxd CSI driver executable accepts following set of command line arguments:

|Name|Default value|Description|
|--|--|--|
|`--endpoint`|`unix:///var/lib/kubelet/plugins/eosxd.csi.cern.ch/csi.sock`|(string value) CSI endpoint. EOSxd CSI will create a UNIX socket at this location.|
|`--drivername`|`eosxd.csi.cern.ch`|(string value) Name of the driver that is used to link PersistentVolume objects to EOSxd CSI driver.|
|`--nodeid`|_none, required_|(string value) Unique identifier of the node on which the EOSxd CSI node plugin pod is running. Should be set to the value of `Pod.spec.nodeName`.|
|`--automount-startup-timeout`|_10_|number of seconds to wait for automount daemon to start up before exiting. `0` means no timeout.|
|`--role`|_none, required_|Enable driver service role (comma-separated list or repeated `--role` flags). Allowed values are: `identity`, `node`, `controller`.|
|`--version`|_false_|(boolean value) Print driver version and exit.|

## automount-runner command line arguments

|Name|Default value|Description|
|--|--|--|
|`--unmount-timeout`|_-1_|number of seconds of idle time after which an autofs-managed EOSxd mount will be unmounted. `0` means never unmount.|
|`--version`|_false_|(boolean value) Print driver version and exit.|

## mount-reconciler command line arguments

|Name|Default value|Description|
|--|--|--|
|`--frequency`|`30s`|(Golang [`time.Duration`](https://pkg.go.dev/time#Duration) value, see [`time.ParseDuration`](https://pkg.go.dev/time#ParseDuration) for accepted values) How often to check mounts in `/eos`|
|`--version`|_false_|(boolean value) Print driver version and exit.|
