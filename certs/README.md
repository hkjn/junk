# certs

Some experiments around dual auth x509 certs.

Generate certificates for CA, server and client with:

```
./make_certs.sh
```

Run the server locally with:
```
go run server/server.go
```

In another terminal, run the client with with:
```
go run client/client.go
```
