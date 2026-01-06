package config

import (
	"os"

	"github.com/shft1/grpc-notes/observability/logger"
)

type AppEnv struct {
	PortGRPC string
}

func SetupAppEnv(log logger.Logger) *AppEnv {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		log.Info("port variable is not found; the default value is 50051")
		port = "50051"
	}
	return &AppEnv{
		PortGRPC: port,
	}
}
