#!/bin/bash

set -euo pipefail

HOST=${HOST:-""}
[ "$HOST" ] || {
	echo "FATAL: No HOST specified." >&2
	exit 1
}
p="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo "Syncing local directory $p to $HOST:src/. Remove the .sync_active file to stop the sync."
PAUSE=${PAUSE:-2s}
touch .sync_active
while true; do
	[ -e .sync_active ] && {
		# --stats
		# time
		rsync -az --progress --stats --delete --exclude=.git/ --exclude=.sync_active --exclude build/ --exclude .gradle/ $p $HOST:src/
		sleep $PAUSE
	} || break
done
echo "No .sync_active file. Exiting."
