#!/bin/bash
#
# Creates CoreOS cluster on GCE.
#
PROJECT_NAME=henrik-jonsson # GCE project, with billing enabled
ZONE=europe-west1-b         # GCE zone
MACHINE_TYPE=f1-micro       # GCE machine type
#MACHINE_TYPE=n1-standard-1
# All valid images can be shown with:
# gcloud compute images list --project coreos-cloud
gcloud compute --project ${PROJECT_NAME} instances create core1 core2 \
		--image https://www.googleapis.com/compute/v1/projects/coreos-cloud/global/images/coreos-stable-522-4-0-v20150108 \
		--zone ${ZONE} --machine-type ${MACHINE_TYPE} --metadata-from-file user-data=cloud-config.yaml

