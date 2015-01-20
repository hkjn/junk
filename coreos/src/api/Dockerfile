FROM golang
MAINTAINER Henrik Jonsson <me@hkjn.me>

RUN go get hkjn.me/junk/coreos/src/api/cmd/apiserver
RUN go install hkjn.me/junk/coreos/src/api/cmd/apiserver

ENTRYPOINT ["/go/bin/apiserver", "-alsologtostderr", "-api_version", "api.v1d8d6af"]

EXPOSE 9100