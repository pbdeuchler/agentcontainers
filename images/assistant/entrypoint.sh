#!/usr/bin/env bash
set -euo pipefail

ln -sf /mnt/state/__store.db /root/.claude/__store.db
ln -sf /mnt/state/projects /root/.claude/projects

exec /root/.local/bin/shim "$@"
