# prototest

## Regenerating protobuf files

```
$ protoc -I report/ report/report.proto --go_out=plugins=grpc:report
```

## Building

```
$ CGO_ENABLED=0 go build -o greeter_client ./client/
$ CGO_ENABLED=0 go build -o greeter_server ./server/
```

Set `GOOS=linux GOARCH=arm` in environment to build towards armv7l.
