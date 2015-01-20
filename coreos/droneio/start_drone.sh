#!/bin/bash

docker run -d --name="drone-test" \
  --env DRONE_GITHUB_CLIENT=${DRONE_GITHUB_CLIENT} \
  --env DRONE_GITHUB_SECRET=${DRONE_GITHUB_SECRET} \
  -p 8080:8080 \
	-v /var/lib/drone/ \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v /home/core/droneio/drone.sqlite:/var/lib/drone/drone.sqlite \
	hkjn/drone:latest
