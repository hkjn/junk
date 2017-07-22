#!/bin/bash
docker run \
        -e RSA_PRIVATE_KEY="$(cat ~/.abuild/mykey.rsa)" \
        -e RSA_PRIVATE_KEY_NAME=mykey.rsa \
        -v "$PWD:/home/builder/package" \
        -v "$HOME/.abuild/packages:/packages" \
        -v "$HOME/.abuild/mykey.rsa.pub:/etc/apk/keys/mykey.rsa.pub" \
        hkjn/alpine-build:$(uname -m)

