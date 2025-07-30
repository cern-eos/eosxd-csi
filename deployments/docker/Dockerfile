FROM gitlab-registry.cern.ch/linuxsupport/alma9-base:20240501-1

# Setup EOS repository.
COPY deployments/docker/eos9al-stable.repo /etc/yum.repos.d/eos9al-stable.repo

# Install the EOS client, autofs-related packages and remove autofs
# configuration files, as they break certain mounts.
#
# eosxd-csi helm deployment comes packaged with it's own autofs configuration.
RUN dnf install -y \
        dnf-plugins-core \
        epel-release && \
    dnf config-manager --set-enabled crb && \
    dnf install -y --setopt=skip_missing_names_on_install=False \
        autofs \
        cern-eos-autofs \
        cern-eos-autofs-squashfs \
        eos-fusex-5.3.14 && \
	dnf clean all && \
    rm -rf \
        /etc/auto.eos.main.misc    \
        /etc/auto.eos.main.project \
        /etc/auto.eos.main.user \
        /usr/local/share/man \
        /usr/share/man \
        /var/cache/dnf

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
