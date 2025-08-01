FROM ubuntu:22.04

# Install dependencies
# RUN apt-get update && apt-get install -y \
#   curl \
#   git \
#   python3 \
#   python3-pip \
#   && rm -rf /var/lib/apt/lists/*
#
# # Install Claude Code
# RUN curl -fsSL http://claude.ai/install.sh | bash
#
#
# ENV PATH="/root/.local/bin:${PATH}"
#
# CMD ["/bin/bash"]

# RUN useradd -m -u 1000 node && \
#   usermod -aG sudo node

# Install basic development tools and iptables/ipset
RUN apt update && apt install -y --no-install-recommends \
  less \
  git \
  procps \
  sudo \
  fzf \
  zsh \
  man-db \
  unzip \
  gnupg2 \
  gh \
  iptables \
  ipset \
  iproute2 \
  dnsutils \
  aggregate \
  jq \
  nano \
  vim \
  wget \
  curl \
  ca-certificates \
  && apt clean && rm -rf /var/lib/apt/lists/*

# Persist bash history.
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  && mkdir /commandhistory \
  && touch /commandhistory/.bash_history
# && chown node:node /commandhistory

# Set `DEVCONTAINER` environment variable to help with orientation
ENV DEVCONTAINER=true

# # Create workspace and config directories and set permissions
# RUN mkdir -p /workspace && \
#   chown -R node /workspace

WORKDIR /workspace

RUN update-ca-certificates -f

ARG GIT_DELTA_VERSION=0.18.2
RUN ARCH=$(dpkg --print-architecture) && \
  wget "https://github.com/dandavison/delta/releases/download/${GIT_DELTA_VERSION}/git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb" && \
  sudo dpkg -i "git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb" && \
  rm "git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb"

# Set up non-root user
# USER node

# Set the default shell to zsh rather than sh
ENV SHELL=/bin/zsh

# Set the default editor and visual
ENV EDITOR nano
ENV VISUAL nano

# Default powerline10k theme
ARG ZSH_IN_DOCKER_VERSION=1.2.0
RUN sh -c "$(wget -O- https://github.com/deluan/zsh-in-docker/releases/download/v${ZSH_IN_DOCKER_VERSION}/zsh-in-docker.sh)" -- \
  -p git \
  -p fzf \
  -a "source /usr/share/doc/fzf/examples/key-bindings.zsh" \
  -a "source /usr/share/doc/fzf/examples/completion.zsh" \
  -a "export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  -x

# # Install Claude Code
RUN curl -fsSL http://claude.ai/install.sh | bash

# # Copy config directories
COPY --chown=1000:1000 .assistant.json /root/.claude.json
COPY --chown=1000:1000 assistant/ /root/.claude

# Copy and set up firewall script
# COPY init-firewall.sh /usr/local/bin/
# USER root
# RUN chmod +x /usr/local/bin/init-firewall.sh && \
#   echo "node ALL=(root) NOPASSWD: /usr/local/bin/init-firewall.sh" > /etc/sudoers.d/node-firewall && \
#   chmod 0440 /etc/sudoers.d/node-firewall
# USER node

ENV PATH="/root/.local/bin:${PATH}"

CMD ["/bin/bash"]
