#!/bin/bash
#
# Submits and starts fleet services.
#
for s in services/*.service; do
		echo "Starting ${s}.."
		fleetctl submit ${s}
		fleetctl start ${s}
done
