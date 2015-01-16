#!/bin/bash
#
# Rebuilds the docker image for DB and pushes it to registry.
#
docker build -t hkjn/coreosdb:latest .
docker push hkjn/coreosdb:latest
