# go-serve


##### Homework 1

Client-Server app that uses goroutines. Every goroutines simulates 1 client.

Connection client-server is using gRPC.

In order to generate go files used in project by the .proto file use the following command:
```go
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto-files/route_guide.proto

```