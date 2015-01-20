FROM golang
MAINTAINER Henrik Jonsson <me@hkjn.me>

RUN go get hkjn.me/junk/coreos/src/web
RUN go install hkjn.me/junk/coreos/src/web

ENTRYPOINT ["/go/bin/web", "-alsologtostderr", "-web_version", "web.v1d8d6af"]

EXPOSE 9000