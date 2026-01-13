package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/shft1/grpc-notes/internal/client/config"
	noteGW "github.com/shft1/grpc-notes/internal/client/gateway/notes/v1"
	noteHand "github.com/shft1/grpc-notes/internal/client/handler/notes/v1"
	noteRoute "github.com/shft1/grpc-notes/internal/client/router/notes/v1"
	"github.com/shft1/grpc-notes/internal/client/server"
	"github.com/shft1/grpc-notes/observability/logger"
	pb "github.com/shft1/grpc-notes/pkg/api/notes/v1"
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

	conn, err := grpc.NewClient(net.JoinHostPort(cfg.HostGRPC, cfg.PortGRPC), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zlog.Error("failed to create connection with gRPC-server", logger.NewField("error", err))
		return
	}
	defer conn.Close()

	notePB := pb.NewNoteAPIClient(conn)
	noteGW := noteGW.NewNoteGateway(zlog, notePB)
	noteHand := noteHand.NewNoteHandler(zlog, noteGW)

	router := chi.NewRouter()

	noteRoute := noteRoute.NewNoteRouter(router, noteHand)
	noteRoute.SetupRoutesV1()

	srv := server.NewHTTPServer(zlog, router, cfg)
	srv.StartGracefully(sysCtx)
}
