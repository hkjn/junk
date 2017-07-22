#!/bin/bash

set -euo pipefail

HOST=${1:-""}
[[ "$HOST" ]] || {
	echo "FATAL: No HOST specified." >&2
	exit 1
}
DIR="$(pwd)"
echo "Syncing local directory $DIR to $HOST:src/hkjn.me. Remove the .sync_active file to stop the sync."
PAUSE=${PAUSE:-2s}
touch .sync_active
while true; do
	[[ -e .sync_active ]] || break
	rsync -az --exclude=.git/ --exclude=.sync_active $DIR $HOST:src/hkjn.me/
	sleep $PAUSE
done
echo "No .sync_active file. Exiting."
