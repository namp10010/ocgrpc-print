package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "../proto"
	grpc "google.golang.org/grpc"
    "go.opencensus.io/examples/exporter"
    "go.opencensus.io/plugin/ocgrpc"
    "go.opencensus.io/stats/view"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
    // register exporter to print to stdout
    view.RegisterExporter(&exporter.PrintExporter{})

    // register view to collect gRPC client stats
    if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
        log.Fatal(err)
    }

    // establish a connection to the grpc server with stats handler being ocgrpc handler
	grpcConn, err := grpc.Dial(address, grpc.WithStatsHandler(&ocgrpc.ClientHandler{}), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer grpcConn.Close()
	greeterClient := pb.NewGreeterClient(grpcConn)

	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

    view.SetReportingPeriod(time.Second)

    for {
        _, err := greeterClient.Greet(context.Background(), &pb.HelloRequest{Name: name})
        if err != nil {
            log.Printf("could not great: %v", err)
        } 

        time.Sleep(2 * time.Second)
    }
}
