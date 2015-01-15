#!/bin/bash
#
# Rebuilds the docker image and pushes it to registry.
#
docker build -t hkjn/junkapp:v1 .
docker push hkjn/junkapp:v1
