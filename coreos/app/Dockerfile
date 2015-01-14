FROM golang
MAINTAINER Henrik Jonsson <me@hkjn.me>

ADD . /go/src/hkjn.me/junk/coreos/app

# RUN go build app.go
RUN go install hkjn.me/junk/coreos/app

ENTRYPOINT /go/bin/app

EXPOSE 8080