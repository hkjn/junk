FROM golang
MAINTAINER Henrik Jonsson <me@hkjn.me>

RUN go get hkjn.me/junk/coreos/src/web
RUN go install hkjn.me/junk/coreos/src/web

ENTRYPOINT /go/bin/web

EXPOSE 9000