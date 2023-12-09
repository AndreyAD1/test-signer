package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AndreyAD1/test-signer/internal/configuration"
)

type Server struct {
	shutdownFuncs   []func()
}

func NewServer(ctx context.Context, config configuration.ServerConfig) (*Server, error) {
	return &Server{}, nil
}

func (s *Server) Run(ctx context.Context) error {
	return nil
}

func (s *Server) Shutdown(timeout time.Duration) {
	// set the timeout to prevent a system hang
	timeoutFunc := time.AfterFunc(timeout, func() {
		logMsg := fmt.Sprintf(
			"timeout %v has been elapsed, force exit",
			timeout.Seconds(),
		)
		log.Fatal(logMsg)
	})
	defer timeoutFunc.Stop()
	for _, f := range s.shutdownFuncs {
		f()
	}
}