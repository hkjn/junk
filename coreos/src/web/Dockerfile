FROM golang
MAINTAINER Henrik Jonsson <me@hkjn.me>

ADD . /go/src/hkjn.me/junk/coreos/src/web

RUN go install hkjn.me/junk/coreos/src/web

ENTRYPOINT /go/bin/web

EXPOSE 9000