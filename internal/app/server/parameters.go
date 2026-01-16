package server

import (
	"time"

	"google.golang.org/grpc/keepalive"
)

type option func(*parameters)

type parameters struct {
	port string
	keepalive.ServerParameters
}

func setupParameters(opts ...option) *parameters {
	parameters := new(parameters)
	for _, opt := range opts {
		opt(parameters)
	}
	return parameters
}

func WithPort(port string) option {
	return func(parameters *parameters) {
		if port == "" {
			port = "8080"
		}
		parameters.port = port
	}
}

func WithMaxConnectionIdle(tm time.Duration) option {
	return func(parameters *parameters) {
		if tm == 0 {
			return
		}
		parameters.MaxConnectionIdle = tm
	}
}

func WithMaxConnectionAge(tm time.Duration) option {
	return func(parameters *parameters) {
		if tm == 0 {
			return
		}
		parameters.MaxConnectionAge = tm
	}
}

func WithMaxConnectionAgeGrace(tm time.Duration) option {
	return func(parameters *parameters) {
		if tm == 0 {
			return
		}
		parameters.MaxConnectionAgeGrace = tm
	}
}

func WithTime(tm time.Duration) option {
	return func(parameters *parameters) {
		if tm == 0 {
			return
		}
		parameters.Time = tm
	}
}

func WithTimeout(tm time.Duration) option {
	return func(parameters *parameters) {
		if tm == 0 {
			return
		}
		parameters.Timeout = tm
	}
}
