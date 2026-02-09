#!/bin/sh
set -eu

KEY_PATH="${SSH_HOST_KEY_PATH:-/data/host_ed25519}"
KEY_DIR="$(dirname "$KEY_PATH")"

mkdir -p "$KEY_DIR"

if [ ! -f "$KEY_PATH" ]; then
	ssh-keygen -t ed25519 -f "$KEY_PATH" -N ""
fi

exec "$@"
