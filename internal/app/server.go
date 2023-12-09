package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	h "github.com/AndreyAD1/test-signer/internal/app/handlers"
	r "github.com/AndreyAD1/test-signer/internal/app/infrastructure/repositories"
	"github.com/AndreyAD1/test-signer/internal/app/services"
	"github.com/AndreyAD1/test-signer/internal/configuration"
)

type Server struct {
	shutdownFuncs []func()
	httpServer    *http.Server
}

var defaultTimeout = 5

func NewServer(ctx context.Context, config configuration.ServerConfig) (*Server, error) {
	signatureRepo, err := r.NewSignatureCollection(ctx, config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	signatureSvc := services.NewSignatureSvc(signatureRepo, config.PublicKeyFile)
	handlers := h.HandlerContainer{
		ApiSecret:    config.APISecret,
		SignatureSvc: signatureSvc,
		Timeout:      time.Duration(defaultTimeout),
	}

	srvMux := http.NewServeMux()
	srvMux.HandleFunc("/api/v1/sign", handlers.SignAnswersHandler())
	srvMux.HandleFunc("/api/v1/verify", handlers.VerifySignatureHandler())
	httpServer := http.Server{
		Addr:    config.ServerAddress,
		Handler: srvMux,
	}
	return &Server{httpServer: &httpServer}, nil
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

func (s *Server) Run(ctx context.Context) error {
	ctx, cancelCtx := context.WithCancel(ctx)
	idleConnectionsClosed := make(chan struct{})

	go func() {
		signalCh := make(chan os.Signal, 4)
		signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		select {
		case sig := <-signalCh:
			log.Printf("receive an OS signal '%v'", sig)
		case <-ctx.Done():
			log.Printf("start shutdown because of context")
		}

		shutdownCtx, shutdownCtxCancel := context.WithTimeout(ctx, 5*time.Second)
		defer shutdownCtxCancel()
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("a server shutdown error: %v", err)
		}
		close(idleConnectionsClosed)
		cancelCtx()
	}()

	go func() {
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("an HTTP server runtime error: %v", err)
		}
		cancelCtx()
	}()

	<-idleConnectionsClosed
	return nil
}
