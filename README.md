# grpc-pie

simple tool to make requests to the grpc services (heavy in development)

## How to use
Assume we have a grpc service described like
```java
syntax = "proto3";

package pb;

service Documents {
    rpc FindByID (FindByIDRequest) returns (FindResponse);
}

message FindByIDRequest {
    uint64 ID = 1;
}

message FindResponse {
    Document document = 1;
}

message Document {
    uint64 ID = 1;
    string name = 2;
}

```
running on port 30000

so wie can use grppie to request documents

> grpcpie -proto doc.proto localhost:50000/Documents/FindByID ID=1
