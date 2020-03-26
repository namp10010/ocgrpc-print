package main

import (
    "context"
    "log"
    "math/rand"
    "net"
    "net/http"
    "time"

    pb "../proto"
    grpc "google.golang.org/grpc"
    "go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
)

const (
    port = ":50051"
)

type greeterServer struct {
}

func (s *greeterServer) Greet(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
    //comment out the log to clear out the output of the tracing
    //log.Printf("Received request: %v", in.GetName())

    // start the span
    ctx, span := trace.StartSpan(ctx, "my-server-span")
    defer span.End()

    // add some randomness into the response time
    time.Sleep(time.Duration(rand.Float64() * float64(time.Second)))

    return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
    // start zpages for local view
    go func() {
        mux := http.NewServeMux()
        zpages.Handle(mux, "/debug")
        // I don't understand the use of log.Fatal here
        log.Fatal(http.ListenAndServe("127.0.0.1:8081", mux))
    }()

    // register print exporter to print to the stdout
    view.RegisterExporter(&exporter.PrintExporter{})

    // register the view to collect server request count
    if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
        log.Fatal(err)
    }

    // create a tcp listener on given port
    tcpListener, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    // create new grpc server with stats handler being ocgrpc
    grpcServer := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))

    // register the greeter server to serve grpc server
    pb.RegisterGreeterServer(grpcServer, &greeterServer{})

    // grpc server serve at tcp listener
    if err := grpcServer.Serve(tcpListener); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
