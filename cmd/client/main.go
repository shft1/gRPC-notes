package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	openapi "github.com/shft1/grpc-notes/docs/api/notes/v1"
	"github.com/shft1/grpc-notes/internal/client/config"
	noteGW "github.com/shft1/grpc-notes/internal/client/gateway/notes/v1"
	noteHand "github.com/shft1/grpc-notes/internal/client/handler/notes/v1"
	"github.com/shft1/grpc-notes/internal/client/interceptor/stream"
	"github.com/shft1/grpc-notes/internal/client/middleware/cors"
	noteRoute "github.com/shft1/grpc-notes/internal/client/router/notes/v1"
	"github.com/shft1/grpc-notes/internal/client/server"
	"github.com/shft1/grpc-notes/observability/logger"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
	"github.com/shft1/grpc-notes/static/swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	sysCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	zlog, err := logger.NewZapAdapter()
	if err != nil {
		log.Printf("failed to create logger: %v", err)
		return
	}
	defer zlog.Sync()

	if err := godotenv.Load(); err != nil {
		zlog.Warn(".env file not found")
	}
	cfg := config.SetupClientEnv(zlog)

	logStreamInter := stream.NewLoggerInterceptor(zlog)

	conn, err := grpc.NewClient(
		net.JoinHostPort(cfg.HostGRPC, cfg.PortGRPC),
		grpc.WithChainStreamInterceptor(logStreamInter),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		zlog.Error("failed to create connection with gRPC-server", logger.NewField("error", err))
		return
	}
	defer conn.Close()

	notePB := pb.NewNoteAPIClient(conn)

	noteGW := noteGW.NewNoteGateway(zlog, notePB)
	noteHand := noteHand.NewNoteHandler(sysCtx, zlog, noteGW)

	router := chi.NewRouter()
	mainRouter := noteRoute.NewNoteRouter(router, noteHand)

	mainRouter.SetupRoutesV1()

	gwRouter := runtime.NewServeMux()
	if err = pb.RegisterNoteAPIHandlerFromEndpoint(
		sysCtx,
		gwRouter,
		net.JoinHostPort(cfg.HostGRPC, cfg.PortGRPC),
		[]grpc.DialOption{
			grpc.WithChainStreamInterceptor(logStreamInter),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	); err != nil {
		zlog.Error("failed to register endpoints to gateway router", logger.NewField("error", err))
		return
	}
	corsMW := cors.NewCORSMiddleware(zlog)

	mainRouter.SetupGenRoutesV1(gwRouter, corsMW, swagger.Content, openapi.Content)

	srv := server.NewHTTPServer(zlog, router, cfg.Host, cfg.Port)
	srv.StartGracefully(sysCtx)
}
