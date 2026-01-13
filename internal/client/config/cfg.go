package config

import (
	"os"

	"github.com/shft1/grpc-notes/observability/logger"
)

type ClientEnv struct {
	Host     string
	Port     string
	HostGRPC string
	PortGRPC string
}

func SetupClientEnv(log logger.Logger) *ClientEnv {
	host, port := os.Getenv("CLIENT_HOST"), os.Getenv("CLIENT_PORT")
	portGRPC, hostGRPC := os.Getenv("GRPC_PORT"), os.Getenv("GRPC_HOST")
	switch "" {
	case host:
		log.Info("host variable is not found; the default value is localhost")
		port = "localhost"
	case port:
		log.Info("port variable is not found; the default value is 8080")
		port = "8080"
	case portGRPC:
		log.Info("portGRPC variable is not found; the default value is 50051")
		portGRPC = "50051"
	case hostGRPC:
		log.Info("hostGRPC variable is not found; the default value is localhost")
		hostGRPC = "localhost"
	}
	return &ClientEnv{
		Port:     port,
		HostGRPC: hostGRPC,
		PortGRPC: portGRPC,
	}
}
