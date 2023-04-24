FROM registry.cern.ch/magnum/eosd:5.1.5-1.1

LABEL description="EOSxd CSI Plugin"

RUN yum install -y autofs cern-eos-autofs     \
    && yum clean all && rm -rf /var/cache/yum

# eosxd-csi deployment with Helm should provide
# its own autofs configuration. The files listed below
# break certain mounts, and so need to be removed.
RUN rm /etc/auto.eos.main.misc    \
       /etc/auto.eos.main.project \
       /etc/auto.eos.main.user

ADD bin/csi-driver /csi-driver
ADD bin/automount-runner /automount-runner
