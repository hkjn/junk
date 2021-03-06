#!/bin/bash
#
# Rebuilds the docker image and pushes it to registry.
#
# Note that we specify --no-cache, since Docker otherwise would cache
# the first go packages retrieved by "go get". This is fine since the
# image is so small, but if it was larger we'd need to bust the cache
# at the appropriate point in the Dockerfile by using a unique command
# instead.
#
VERSION=api.v$(git rev-parse --short HEAD)
sed -i "s/-web_version web.v.*\b/-web_version ${VERSION}/" Dockerfile
docker build --no-cache -t hkjn/coreosweb:latest .
docker push hkjn/coreosweb:latest
