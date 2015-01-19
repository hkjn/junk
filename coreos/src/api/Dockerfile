FROM golang
MAINTAINER Henrik Jonsson <me@hkjn.me>

RUN go get hkjn.me/junk/coreos/src/api/cmd/apiserver
RUN go install hkjn.me/junk/coreos/src/api/cmd/apiserver

ENTRYPOINT /go/bin/apiserver -alsologtostderr -build_version api.vd56f8c2

EXPOSE 9100