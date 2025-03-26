# How to use EOSxd CSI driver in Kubernetes

`eosxd-csi` supports mounting EOS instances and exposing them inside your Pods as a single PVC.

## Getting Started

Ensure that you have the appropriate `StorageClass` configured in your cluster this can be done either be enabling `commonStorageClass.enabled` on installation or by running the below.

```bash
cat << EOF | kubectl apply -f -
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: eos
provisioner: eosxd.csi.cern.ch
EOF
```

Provision a `PVC` using this `StorageClass` and mount this to the `Pod`'s which require access to `EOS`.

```bash
cat << EOF | kubectl apply -f -
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: eos
spec:
  accessModes:
  - ReadWriteMany
  resources:
    requests:
      # Volume size value has no effect and is ignored
      # by the driver, but must be non-zero.
      storage: 1
  storageClassName: eos

EOF
```

Please note that when mounting the PVC, the Pod’s volumeMounts spec of the volume must set `mountPropagation: HostToContainer`.

```yaml
...
    volumeMounts:
    - name: eos
      mountPath: /eos
      mountPropagation: HostToContainer
...
```

`eosxd-csi` exposes the autofs-EOS root in `/var/eos` on the host node, whilst mounting via `PVC` is preferred you can use this `hostPath`.

```yaml
...
  containers:
  - volumeMounts:
    - name: eos
      mountPath: /eos
      mountPropagation: HostToContainer
...
  volumes:
  - name: eos
    hostPath:
      path: /var/eos
      type: Directory
...
```

## Authentication

Access to EOS requires authentication. Popular options are:

* Kerberos tickets
* OAuth2 tokens

### Kerberos

Authentication using Kerberos tickets requires that your Pod to come with Kerberos tools installed.

Steps:

1. Create a Pod mounting the EOS PVC.
2. Inside the Pod, use `kinit <CERN login name>@CERN.CH` to create a Kerberos ticket.
Example:

```bash
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: eos-kerberos
spec:
  containers:
  - name: my-container
    image: registry.cern.ch/docker.io/cern/alma9-base:latest
    imagePullPolicy: IfNotPresent
    command: ["sleep", "inf"]
    volumeMounts:
    - name: eos
      mountPath: /eos
      mountPropagation: HostToContainer
  volumes:
  - name: eos
    persistentVolumeClaim:
      claimName: eos
EOF

# This example exec's into the eos-kerberos Pod, creates
# a Kerberos ticket for user "rvasek" using the kinit command,
# and lists the home directory.
$ kubectl exec -it eos-kerberos -- bash
[root@eos-kerberos /]# kinit rvasek@CERN.CH
Password for rvasek@CERN.CH: <CERN login password>
[root@eos-kerberos /]# cd /eos/home-r/rvasek
[root@eos-kerberos rvasek]# ls -l
total 24
drwxr-xr-x. 2 110701 2763 4096 Sep  9  2022 Documents
drwxr-xr-x. 2 110701 2763 4096 Feb  2  2020 Music
drwxr-xr-x. 2 110701 2763 4096 Feb  2  2020 Pictures
drwxr-xr-x. 2 110701 2763 4096 May 30 14:59 SWAN_projects
drwxr-xr-x. 2 110701 2763 4096 Feb  2  2020 Videos
...
```

### OAuth2

At CERN, we maintain a `oauth2-refresh-controller` that allows the use of an OIDC application for authentication.

For more details please see the [gitlab repo](https://gitlab.cern.ch/kubernetes/security/oauth2-refresh-controller).

We detail the setup process for using this controller at CERN as an illustrative end-to-end example, but naturally this will differ according to your organisation.

To ensure you have a valid access token at CERN, it must satisfy the following:

* Be created with the openid scope (at least) and an EOS-enabled client ID.
* Stored in `/tmp/oauthtk_<Container's UID>`*.
* Group and ownership set to match the container’s GID and UID*.
* Have file permissions set to 0400.
* Have it's file contents in the format of oauth2:<OAuth2 access token>:auth.cern.ch/auth/realms/cern/protocol/openid-connect/userinfo.

*: or whatever the effective UID and GID are set to in the process accessing the EOS storage)

The example below then creates a `Secret` with an `Oauth2` JWT and a `ConfigMap` with the EOS-specific token format. The `oauth2-refresh-controller` annotation on the pod is then used to inject and continually refresh your EOS token.
```bash
cat << EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: my-oauth2-token
  annotations:
    oauth2-refresh-controller.cern.ch/is-token: "true"
stringData:
  oauth2: "<OAuth2 JWT>"
  clientID: "<OIDC Client ID>"
  clientSecret: "<OIDC Client password>"

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: oauth2-token-eos-template
data:
  # This defines what the token file will contain.
  # The "$(ACCESS_TOKEN)" variable is expanded into
  # the actual value by oauth2-refresh-controller.
  template: "oauth2:$(ACCESS_TOKEN):auth.cern.ch/auth/realms/cern/protocol/openid-connect/userinfo"

---
apiVersion: v1
kind: Pod
metadata:
  name: eos-oauth2
  annotations:
    oauth2-refresh-controller.cern.ch/to-inject: |
      [
        {
          "secretName": "my-oauth2-token",
          "container": "my-container",
          "templateConfigMapName": "oauth2-token-eos-template",
          "owner": 0
        }
      ]      
spec:
  containers:
  - name: my-container
    image: registry.cern.ch/docker.io/cern/alma9-base:latest
    imagePullPolicy: IfNotPresent
    command: ["sleep", "inf"]
    volumeMounts:
    - name: eos
      mountPath: /eos
      mountPropagation: HostToContainer
  volumes:
  - name: eos
    persistentVolumeClaim:
      claimName: eos
EOF

# This example exec's into the eos-oauth2 Pod and lists
# the home directory of the user "rvasek". The OAuth2 access
# token is injected into the container and renewed automatically.

$ kubectl exec -it eos-oauth2 -- bash
Defaulted container "my-container" out of: my-container, oauth2-refresh-bootstrap-0 (init)
[root@eos-oauth2 rvasek]# cd /eos/home-r/rvasek
[root@eos-oauth2 rvasek]# ls -l
total 24
drwxr-xr-x. 2 110701 2763 4096 Sep  9  2022 Documents
drwxr-xr-x. 2 110701 2763 4096 Feb  2  2020 Music
drwxr-xr-x. 2 110701 2763 4096 Feb  2  2020 Pictures
drwxr-xr-x. 2 110701 2763 4096 May 30 14:59 SWAN_projects
drwxr-xr-x. 2 110701 2763 4096 Feb  2  2020 Videos
...
```


For more information on getting started please refer to examples in [../example/](../example/) or CERN's public documentation for [eosxd-csi](https://kubernetes.docs.cern.ch/docs/storage/eos/).

