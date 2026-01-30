package config

import (
	"os"
	"strconv"
	"time"

	"github.com/shft1/grpc-notes/observability/logger"
)

type ServerEnv struct {
	Port                  string
	MaxConnectionIdle     time.Duration
	MaxConnectionAge      time.Duration
	MaxConnectionAgeGrace time.Duration
	Time                  time.Duration
	Timeout               time.Duration
	Capacity              int
}

func SetupServerEnv(log logger.Logger) *ServerEnv {
	port := os.Getenv("GRPC_PORT")
	idle, _ := time.ParseDuration(os.Getenv("MAX_CONNECTION_IDLE"))
	age, _ := time.ParseDuration(os.Getenv("MAX_CONNECTION_AGE"))
	grace, _ := time.ParseDuration(os.Getenv("MAX_CONNECTION_AGE_GRACE"))
	tm, _ := time.ParseDuration(os.Getenv("TIME"))
	tmout, _ := time.ParseDuration(os.Getenv("TIMEOUT"))
	cap, err := strconv.Atoi(os.Getenv("CAPACITY"))
	if err != nil {
		cap = 3
		log.Info("set default value on CAPACITY", logger.NewField("capacity", cap))
	}
	return &ServerEnv{
		Port:                  port,
		MaxConnectionIdle:     idle,
		MaxConnectionAge:      age,
		MaxConnectionAgeGrace: grace,
		Time:                  tm,
		Timeout:               tmout,
		Capacity:              cap,
	}
}
