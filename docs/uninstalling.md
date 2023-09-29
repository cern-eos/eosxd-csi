## Uninstalling eosxd-csi driver

The nodeplugin Pods store various resources on the node hosts they are running on:
* autofs mount and the respective inner EOS mounts,
* EOS client cache.

By default, the nodeplugin Pod leaves autofs and its respective inner mounts on the node
in `/var/eos`. They may need to be unmounted recursively. To do that, you can set
`AUTOFS_TRY_CLEAN_AT_EXIT` environment variable to `true` in nodeplugin's DaemonSet and restart
the Pods. On the next exit, they will be unmounted.

The EOS client cache is stored by default in `/var/cache/eos`. This directory is not deleted
automatically, and manual intervention is currently needed.
