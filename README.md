# Inspr Core
This is the core for Inspr. It's what makes the cluster run.

## Dependencies
```sh
go get -u k8s.io/client-go@v0.17.2 github.com/googleapis/gnostic@v0.3.1
```

## Protobuffers
Protobuffers can be compiled using this command (from within the folder). Assuming you're on folder `pkg/operator/channel` and want to include `pkg/meta.meta.proto` on your protobuffer file, you can use `protoc` by running the include commands below:

```sh
protoc -I../../../ -I.  --go-grpc_out=. channel.proto
```
