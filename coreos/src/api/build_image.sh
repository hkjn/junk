#!/bin/bash
#
# Rebuilds the docker image and pushes it to registry.
#
docker build -t hkjn/coreosapi:latest .
docker push hkjn/coreosapi:latest
