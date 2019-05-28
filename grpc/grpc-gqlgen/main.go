package main

import (
	"context"
	"fmt"
	rlog "log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/vvakame/til/grpc/grpc-gqlgen/graphqlapi"
	"github.com/vvakame/til/grpc/grpc-gqlgen/grpcapi"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/akutz/memconn"
	"github.com/favclip/ucon"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/vvakame/sdlog/aelog"
	"github.com/vvakame/til/grpc/grpc-gqlgen/echopb"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var echoCli echopb.EchoClient

func main() {

	ctx := context.Background()

	if os.Getenv("GOOGLE_CLOUD_PROJECT") != "" {
		exporter, err := stackdriver.NewExporter(stackdriver.Options{
			ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
		})
		if err != nil {
			rlog.Fatalf("Failed to create stackdriver exporter: %v", err)
		}
		trace.RegisterExporter(exporter)
		view.RegisterExporter(exporter)
		trace.ApplyConfig(trace.Config{
			DefaultSampler: trace.AlwaysSample(),
		})
		defer exporter.Flush()
	}

	ucon.Orthodox()

	ucon.HandleFunc("POST", "/echo", echoHandler)

	config, err := graphqlapi.InitializeGraphQLConfig(context.Background())
	if err != nil {
		rlog.Fatal(err)
	}

	ucon.HandleFunc("GET", "/", handler.Playground("GraphQL playground", "/api/query"))
	ucon.HandleFunc("*", "/api/query", handler.GraphQL(
		graphqlapi.NewExecutableSchema(config),
	))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		rlog.Printf("Defaulting to port %s for HTTP", port)
	}

	rlog.Printf("Listening HTTP on port %s", port)

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: &ochttp.Handler{
			Handler:          ucon.DefaultMux,
			IsPublicEndpoint: true,
			Propagation:      &propagation.HTTPFormat{},
		},
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			rlog.Fatal(err)
		}
	}()

	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		rlog.Fatalf("Failed to register gRPC client views: %v", err)
	}
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		rlog.Fatalf("Failed to register gRPC server views: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))
	echoServer, err := grpcapi.NewEchoServer()
	if err != nil {
		rlog.Fatal(err)
	}
	echopb.RegisterEchoServer(grpcServer, echoServer)
	reflection.Register(grpcServer)

	listener, err := memconn.Listen("memu", "grpc")
	if err != nil {
		rlog.Fatal(err)
	}
	defer listener.Close()

	rlog.Print("Listening gRPC on net.Pipe")

	go func() {
		if err := grpcServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			rlog.Fatal(err)
		}
	}()

	conn, err := grpc.Dial(
		"grpc",
		grpc.WithInsecure(),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return memconn.Dial("memu", addr)
		}),
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
	)
	if err != nil {
		rlog.Fatal(err)
	}

	echoCli = echopb.NewEchoClient(conn)

	mux := runtime.NewServeMux()
	err = echopb.RegisterEchoHandlerClient(ctx, mux, echoCli)
	if err != nil {
		rlog.Fatal(err)
	}
	ucon.HandleFunc("*", "/v1/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	}))

	rlog.Printf("running...")

	// setup graceful shutdown...
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		rlog.Fatalf("graceful shutdown failure: %s", err)
	}
	grpcServer.GracefulStop()
	rlog.Printf("graceful shutdown successfully")
}

func echoHandler(ctx context.Context, req *echopb.SayRequest) (*echopb.SayResponse, error) {
	ctx, span := trace.StartSpan(ctx, "echoHandler")
	defer span.End()
	span.AddAttributes(trace.StringAttribute("messageId", req.MessageId))
	span.AddAttributes(trace.StringAttribute("messageBody", req.MessageBody))

	resp, err := echoCli.Say(ctx, req)
	status := status.Convert(err)
	log.Debugf(ctx, "%#v", status)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
