package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	oczipkin "contrib.go.opencensus.io/exporter/zipkin"
	"github.com/go-kit/kit/log/level"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/muhammadisa/bareksanews/constant"
	ep "github.com/muhammadisa/bareksanews/endpoint"
	"github.com/muhammadisa/bareksanews/gvars"
	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
	"github.com/muhammadisa/bareksanews/repository"
	"github.com/muhammadisa/bareksanews/service"
	"github.com/muhammadisa/bareksanews/transport"
	"github.com/muhammadisa/bareksanews/util/cb"
	"github.com/muhammadisa/bareksanews/util/dbc"
	"github.com/muhammadisa/bareksanews/util/hdr"
	"github.com/muhammadisa/bareksanews/util/lgr"
	"github.com/openzipkin/zipkin-go"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/soheilhy/cmux"
	"go.opencensus.io/trace"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// ServeGRPC serving GRPC server and will be merged using MergeServer function
func ServeGRPC(listener net.Listener, service pb.BareksaNewsServiceServer, serverOptions []grpc.ServerOption) error {
	level.Info(gvars.Log).Log(lgr.LogInfo, "initialize grpc server")

	var grpcServer *grpc.Server
	if len(serverOptions) > 0 {
		grpcServer = grpc.NewServer(serverOptions...)
	} else {
		grpcServer = grpc.NewServer()
	}
	pb.RegisterBareksaNewsServiceServer(grpcServer, service)
	return grpcServer.Serve(listener)
}

// ServeHTTP serving HTTP server and will be merged using MergeServer function
func ServeHTTP(listener net.Listener, service pb.BareksaNewsServiceServer) error {
	level.Info(gvars.Log).Log(lgr.LogInfo, "initialize rest server")

	mux := runtime.NewServeMux()
	err := pb.RegisterBareksaNewsServiceHandlerServer(context.Background(), mux, service)
	if err != nil {
		return err
	}
	srv := &http.Server{Handler: hdr.CORS(mux)}
	return srv.Serve(listener)
}

// MergeServer start ServeGRPC and ServeHTTP concurrently
func MergeServer(service pb.BareksaNewsServiceServer, serverOptions []grpc.ServerOption) {
	port := fmt.Sprintf(":%s", "8010")
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	m := cmux.New(listener)
	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings(
		"content-type", "application/grpc",
	))
	httpListener := m.Match(cmux.HTTP1Fast())

	g := new(errgroup.Group)
	g.Go(func() error { return ServeGRPC(grpcListener, service, serverOptions) })
	g.Go(func() error { return ServeHTTP(httpListener, service) })
	g.Go(func() error { return m.Serve() })

	log.Fatal(g.Wait())
}

func main() {
	gvars.Log = lgr.Create(constant.ServiceName)

	level.Info(gvars.Log).Log(lgr.LogInfo, "service started")

	ctx := context.Background()
	defer ctx.Done()

	reporter := httpreporter.NewReporter("http://localhost:9411/api/v2/spans")
	localEndpoint, _ := zipkin.NewEndpoint(constant.ServiceName, ":0")
	exporter := oczipkin.NewExporter(reporter, localEndpoint)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	trace.RegisterExporter(exporter)
	trcr := trace.DefaultTracer

	var repoConf repository.RepoConf
	{
		repoConf.SQL = dbc.Config{
			Username: "root",
			Password: "root",
			Host:     "localhost",
			Port:     "3306",
			Name:     "news_bareksa",
		}
		repoConf.Cache = dbc.Config{
			Password: "password",
			Host:     "localhost",
			Port:     "6379",
		}
	}

	err := cb.StartHystrix(constant.CircuitBreakerTimeout, constant.ServiceName)
	if err != nil {
		panic(err)
	}

	tagRepo, err := repository.NewRepository(ctx, repoConf, trcr)
	if err != nil {
		panic(err)
	}

	bareksaNewsEp, err := ep.NewBareksaNewsEndpoint(service.NewUsecases(*tagRepo, trcr), gvars.Log)
	if err != nil {
		panic(err)
	}

	MergeServer(transport.NewBareksaNewsServer(bareksaNewsEp), nil)
}
