FROM ubuntu:22.04

# Install basic development tools and iptables/ipset
RUN apt-get update && apt install -y --no-install-recommends \
  less \
  git \
  procps \
  fzf \
  zsh \
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
  && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN curl -sL https://deb.nodesource.com/setup_4.x | bash
RUN apt-get install nodejs

# Persist bash history.
RUN SNIPPET="export PROMPT_COMMAND='history -a' && export HISTFILE=/commandhistory/.bash_history" \
  && mkdir /commandhistory \
  && touch /commandhistory/.bash_history

# Set `DEVCONTAINER` environment variable to help with orientation
ENV DEVCONTAINER=true

WORKDIR /workspace

RUN update-ca-certificates -f

ARG GIT_DELTA_VERSION=0.18.2
RUN ARCH=$(dpkg --print-architecture) && \
  wget "https://github.com/dandavison/delta/releases/download/${GIT_DELTA_VERSION}/git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb" && \
  dpkg -i "git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb" && \
  rm "git-delta_${GIT_DELTA_VERSION}_${ARCH}.deb"

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
RUN curl -fsSL https://claude.ai/install.sh | bash

# # Copy config directories
COPY .claude.json /root/.claude.json
COPY .claude/ /root/.claude

COPY .mcp.json /workspace/.mcp.json

RUN mkdir -p /mnt/state

COPY --from=agentcontainers/shim /usr/src/shim/build/shim /root/.local/bin/shim
COPY --from=agentcontainers/mcp /usr/src/mcp/build/assistant-mcp /root/.local/bin/assistant-mcp

COPY entrypoint.sh /root/.local/bin/entrypoint.sh

ENV PATH="/root/.local/bin:${PATH}"
ENV IS_SANDBOX=1
ENV SYSTEM_PROMPT="You are not Claude Code. You are Clarice, a highly trained and experienced a personal assistant for a married couple. You have a wide array of skills: from everything like managing finances, abstract problem solving, home management, childcare expert, marriage counselor, to planning vacations or helping with home improvement projects. You are also a great cook and can help with meal planning and grocery shopping. You are always ready to help with any task, big or small. You are friendly, efficient, and always ready to help."


ENTRYPOINT ["/root/.local/bin/shim"]
