package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/shft1/grpc-notes/internal/app/config"
	noteHand "github.com/shft1/grpc-notes/internal/app/handler/notes/v1"
	noteRepo "github.com/shft1/grpc-notes/internal/app/repository/notes"
	"github.com/shft1/grpc-notes/internal/app/server"
	noteUcase "github.com/shft1/grpc-notes/internal/app/usecase/notes"
	"github.com/shft1/grpc-notes/observability/logger"

	"github.com/joho/godotenv"
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
	cfg := config.SetupAppEnv(zlog)

	noteRepo := noteRepo.NewNoteRepository(zlog)
	noteUcase := noteUcase.NewNotesUseCase(zlog, noteRepo)
	noteHand := noteHand.NewNoteHandler(zlog, noteUcase)

	srv, lis, err := server.NewServer(zlog, cfg)
	if err != nil {
		zlog.Error("failed to create grpc-server", logger.NewField("error", err))
		return
	}
	srv.RegisterRPC(noteHand)

	srv.StartGracefully(sysCtx, zlog, lis)
}
