FROM gitlab-registry.cern.ch/linuxsupport/alma9-base:20240201-2

# Setup EOS repository.
ADD deployments/docker/eos9al-stable.repo /etc/yum.repos.d/eos9al-stable.repo
RUN yum install -y dnf-plugins-core epel-release \
    && dnf config-manager --set-enabled crb

# Install the EOS client and autofs-related packages.
RUN yum install -y --setopt=skip_missing_names_on_install=False \
        eos-fusex-5.2.22                                        \
        autofs                                                  \
        cern-eos-autofs                                         \
        cern-eos-autofs-squashfs                                \
	&& yum clean all

# eosxd-csi deployment with Helm should provide
# its own autofs configuration. The files listed below
# break certain mounts, and so need to be removed.
RUN rm /etc/auto.eos.main.misc    \
       /etc/auto.eos.main.project \
       /etc/auto.eos.main.user

ARG TARGETARCH
ARG RELEASE
ARG GITREF
ARG CREATED

LABEL description="EOSxd CSI Plugin" \
      org.opencontainers.image.title="eosxd-csi" \
      org.opencontainers.image.description="The EOSxd CSI image contains tools for exposing EOS instances through Container Storage Interface." \
      org.opencontainers.image.authors="EOSxd CSI authors" \
      org.opencontainers.image.vendor="EOSxd CSI authors" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.documentation="https://github.com/cern-eos/eosxd-csi" \
      org.opencontainers.image.source="https://github.com/cern-eos/eosxd-csi" \
      org.opencontainers.image.url="https://github.com/cern-eos/eosxd-csi" \
      org.opencontainers.image.revision=${GITREF} \
      org.opencontainers.image.version=${RELEASE} \
      org.opencontainers.image.created=${CREATED} \
      org.opencontainers.image.ref.name="" \
      org.opencontainers.image.base.digest="" \
      org.opencontainers.image.base.name=""

COPY bin/linux-${TARGETARCH}/csi-driver /csi-driver
COPY bin/linux-${TARGETARCH}/automount-runner /automount-runner
COPY bin/linux-${TARGETARCH}/mount-reconciler /mount-reconciler
