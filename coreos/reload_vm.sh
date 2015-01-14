#!/bin/bash
#
# Reloads metadata and resets VM given as first argument.
#
ZONE=europe-west1-b         # GCE zone
gcloud compute instances add-metadata ${1} --zone ${ZONE} --metadata-from-file user-data=cloud-config.yaml
sleep 10
gcloud compute instances reset --zone ${ZONE} ${1}
