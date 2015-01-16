FROM golang
MAINTAINER Henrik Jonsson <me@hkjn.me>

ADD . /go/src/hkjn.me/junk/coreos/src/api

RUN go get hkjn.me/junk/coreos/src/api/cmd/apiserver
RUN go install hkjn.me/junk/coreos/src/api/cmd/apiserver

ENTRYPOINT /go/bin/apiserver

EXPOSE 9100