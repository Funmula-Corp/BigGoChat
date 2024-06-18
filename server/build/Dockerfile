FROM ubuntu:jammy-20230308

# Setting bash as our shell, and enabling pipefail option
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

# Some ENV variables
ENV PATH="/mattermost/bin:${PATH}"

# Build Arguments
ARG PUID=2000
ARG PGID=2000
# MM_PACKAGE build arguments controls which version of mattermost to install, defaults to latest stable enterprise
# i.e. https://releases.mattermost.com/9.7.1/mattermost-9.7.1-linux-amd64.tar.gz
# ARG MM_PACKAGE="https://latest.mattermost.com/mattermost-enterprise-linux"

# # Install needed packages and indirect dependencies
RUN apt-get update \
  && DEBIAN_FRONTEND=noninteractive apt-get install --no-install-recommends -y \
  ca-certificates \
  curl \
  mime-support \
  unrtf \
  wv \
  poppler-utils \
  tidy \
  tzdata \
  net-tools \
  && rm -rf /var/lib/apt/lists/*

# Set mattermost group/user
RUN mkdir -p /mattermost/data /mattermost/plugins /mattermost/client/plugins \
  && addgroup -gid ${PGID} mattermost \
  && adduser -q --disabled-password --uid ${PUID} --gid ${PGID} --gecos "" --home /mattermost mattermost

# get mattermost from make package
COPY packages/dist/mattermost-team-linux-amd64.tar.gz .
RUN tar zxvf mattermost-team-linux-amd64.tar.gz && rm mattermost-team-linux-amd64.tar.gz \
  && chown -R mattermost:mattermost /mattermost /mattermost/data /mattermost/plugins /mattermost/client/plugins

# custom certificate for development ap
COPY Funmula_Internal_Root_CA.crt /usr/local/share/ca-certificates/
RUN update-ca-certificates

# We should refrain from running as privileged user
# USER mattermost
USER root

# Healthcheck to make sure container is ready
HEALTHCHECK --interval=30s --timeout=10s \
  CMD curl -f http://localhost:8065/api/v4/system/ping || exit 1

# Configure entrypoint and command with proper permissions
# COPY --chown=mattermost:mattermost --chmod=765 entrypoint.sh /
ADD entrypoint.sh /
RUN chown mattermost:mattermost /entrypoint.sh
RUN chmod 765 /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
WORKDIR /mattermost
CMD ["mattermost"]

EXPOSE 8065 8067 8074 8075

# Declare volumes for mount point directories
VOLUME ["/mattermost/data", "/mattermost/logs", "/mattermost/config", "/mattermost/plugins", "/mattermost/client/plugins"]
