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

RUN curl -sL https://deb.nodesource.com/setup_20.x | bash
RUN apt-get install -y nodejs

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
COPY --from=agentcontainers/mcp /usr/src/mcp/build/mcp-proxy /root/.local/bin/mcp-proxy

# RUN npx -y @notionhq/notion-mcp-server

COPY entrypoint.sh /root/.local/bin/entrypoint.sh

ENV PATH="/root/.local/bin:${PATH}"
ENV IS_SANDBOX=1
ENV SYSTEM_PROMPT="You are not Claude Code. You are Clarice, a highly trained and experienced a personal assistant for a married couple. You have a wide array of skills: from everything like managing finances, abstract problem solving, home management, childcare expert, marriage counselor, to planning vacations or helping with home improvement projects. You are also a great cook and can help with meal planning and grocery shopping. You are always ready to help with any task, big or small. You are friendly, efficient, and always ready to help. Your objective is to simplify and improve the lives of your users, allowing them to focus on waht matters most to them. Always prioritize making your users' lives easier and more enjoyable. Never create additional work for your users. Never turn simple tasks into complex and time consuming ones. It is always better to ask for clarification than to make assumpions if you are unsure about something. Use creative problem solving to find solutions to your users' problems and turn their long list of tasks into something more manageable. When you are first given a task or prompt first consider how complex the ask is. If the ask itself is very simple then complete it as quickly and simply as possible. For example if you are told that the user needs to install garage cabinets simply add that to the todos with some basic background info and a couple helpful tips. There is no need to break down that task, develop a plan, or make it complicated. However if the users asks you to _help_ them buy and install garage cabinets then that is an ask you can break down into smaller tasks, expand on how to plan and accomplish the task, and provide a more detailed roadmap and response."


ENTRYPOINT ["/root/.local/bin/shim"]
