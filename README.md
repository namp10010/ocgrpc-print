# Example gRPC server and client with OpenCensus

This is a modified version from [gRPC with OpenCensus](https://github.com/census-instrumentation/opencensus-go/tree/master/examples/grpc)

This example uses:

* gRPC to create an RPC server and client.
* The OpenCensus gRPC plugin to instrument the RPC server and client.
* Debugging exporters to print stats and traces to stdout.

First, run the server:

```
$ go run server/main.go
```

Then, run the client:

```
$ go run client/main.go
```

You will see traces and stats exported on the stdout. You can use one of the
[exporters](https://godoc.org/go.opencensus.io/exporter)
to upload collected data to the backend of your choice.

You can also see the z-pages provided from the server:
* Traces: http://localhost:8081/debug/tracez
* RPCs: http://localhost:8081/debug/rpcz
