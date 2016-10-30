#
# Bootstrap Certificate Authorities (infra-etcd-ca, k8s-etcd-ca,
# k8s-apiserver-ca) and generate certs and keys for apiserver,
# kubectl, and etcd.
#
set -euo pipefail

cd caconfig

makeCA() {
  local ca=$1
  echo "Making CA '$ca'.."
  cat ca_csr.json | cfssl gencert -initca - | cfssljson -bare "../.certs/$ca"
}

makeEcdsaCert() {
  local ca=$1
  local cert=$2
  local hostnames="127.0.0.1"
  if [ $# -eq 3 ]; then
    hostnames="$hostnames,$3"
  fi
  echo "Making ECDA cert for '$cert' against CA '$ca' for hostnames '$hostnames'.."
  cat ecdsa_csr.json | cfssl gencert -ca=../.certs/${ca}.pem -ca-key=../.certs/${ca}-key.pem -config=ca_config.json -profile=server -hostname="$hostnames" - | cfssljson -bare "../.certs/$cert"
}

makeRsaCert() {
  local ca=$1
  local cert=$2
  local hostnames="$3"
  echo "Making RSA cert for '$cert' against CA '$ca' for hostnames '$hostnames'.."
  cat rsa_csr.json | cfssl gencert -ca=../.certs/${ca}.pem -ca-key=../.certs/${ca}-key.pem -config=ca_config.json -profile=server -hostname="$hostnames" - | cfssljson -bare "../.certs/$cert"
}

mkdir -p ../.certs

makeCA infra-etcd-ca
makeCA k8s-etcd-ca
makeCA k8s-apiserver-ca

# TODO: Take m1.tf.hkjn.me as var?
makeEcdsaCert k8s-apiserver-ca k8s-apiserver "10.100.0.1,m1.tf.hkjn.me"
makeRsaCert k8s-apiserver-ca service-account "10.100.0.1"
makeEcdsaCert k8s-apiserver-ca apiserver-client
makeEcdsaCert k8s-apiserver-ca kubectl-client

makeEcdsaCert k8s-etcd-ca etcd "m1.tf.hkjn.me"
makeEcdsaCert k8s-etcd-ca etcd-client

makeEcdsaCert infra-etcd-ca etcd "m1.tf.hkjn.me"
makeEcdsaCert infra-etcd-ca etcd-client

chmod 400 ../.certs/{*.pem,*.csr}

echo "All done; CAs and certs are in ../.certs."
