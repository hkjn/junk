#!/bin/bash
#
# Connects to mysql running inside Docker container.
#
CONTAINER_ID=$(docker ps | grep db | awk '{print $1}')
IP=$(docker inspect --format '{{.NetworkSettings.IPAddress}}' $CONTAINER_ID)
mysql -u dbuser -p -h $IP
