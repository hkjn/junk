#!/bin/bash

set -euo pipefail

log() {
  echo "$@" >> /var/log/cloud-init-master.log
}

log "Adding kubernetes key and repos.."
curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
cat <<EOF > /etc/apt/sources.list.d/kubernetes.list
deb http://apt.kubernetes.io/ kubernetes-xenial main
EOF

log "Installing updates.."
apt-get update

log "Installing docker.."
apt-get -y install docker

log "Installing kubernetes.."
apt-get install -y kubelet kubeadm kubectl kubernetes-cni

log "Running kubeadm init.."
kubeadm init --token "${k8s_token}"

log "All done running master.sh."
