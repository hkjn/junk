#!/bin/bash
#
# Rebuilds the docker image and pushes it to registry.
#
docker build -t hkjn/coreosweb:latest .
docker push hkjn/coreosweb:latest
