#!/bin/bash
#
# Create CA, server and client certificates.
#

set -euo pipefail

echo '{"CN":"CA","key":{"algo":"rsa","size":4096}}' | cfssl gencert -initca - | cfssljson -bare ca -
echo '{"signing":{"default":{"expiry":"43800h","usages":["signing","key encipherment","server auth","client auth"]}}}' > ca-config.json
ADDRESS=dev.tf.hkjn.me
NAME=server
echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":4096}}' | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - | cfssljson -bare $NAME
ADDRESS=
NAME=client
echo '{"CN":"'$NAME'","hosts":[""],"key":{"algo":"rsa","size":4096}}' | cfssl gencert -config=ca-config.json -ca=ca.pem -ca-key=ca-key.pem -hostname="$ADDRESS" - | cfssljson -bare $NAME

rm *.csr
rm ca-config.json

openssl x509 -in ca.pem -text -noout
openssl x509 -in server.pem -text -noout
openssl x509 -in client.pem -text -noout
