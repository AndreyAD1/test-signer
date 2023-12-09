package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AndreyAD1/test-signer/internal/app"
	"github.com/AndreyAD1/test-signer/internal/configuration"
	"github.com/caarlos0/env/v9"
	"github.com/spf13/cobra"
)

var (
	apiSecret string
	databaseURL string
	debug bool
	RootCmd = &cobra.Command{
		Use:   "test-signer",
		Short: "The 'Test Signer' service.",
		Long: `The Test signer is a service that accepts a set of answers and 
	questions and signs that the user has finished the " test " at this point in time. 
	The signatures are stored and can later be verified by a different service.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
	panicCounter int
	panicThreshold = 10
)

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.Flags().StringVarP(
		&apiSecret,
		"secret",
		"s",
		"",
		"a secret to manage a JWT token",
	)
	RootCmd.Flags().StringVarP(
		&databaseURL,
		"dburl",
		"u",
		"",
		"a database URL",
	)
	RootCmd.Flags().BoolP(
		"debug",
		"d",
		false,
		"a debug mode",
	)
}

func run() error {
	defer func() {
		p := recover()
		if p == nil {
			return
		}
		panicCounter++
		if panicCounter >= panicThreshold {
			log.Printf("too many panics: %v", panicCounter)
			return
		}
		run()
	}()

	if apiSecret != "" {
		os.Setenv("API_SECRET", apiSecret)
	}
	if databaseURL != "" {
		os.Setenv("DATABASE_URL", databaseURL)
	}
	if debug {
		os.Setenv("DEBUG", "true")
	}
	config := configuration.ServerConfig{}
	err := env.Parse(&config)
	if err != nil {
		return fmt.Errorf("a configuration error: %w", err)
	}

	ctx := context.Background()
	server, err := app.NewServer(ctx, config)
	if err != nil {
		return fmt.Errorf("can not create a new server: %w", err)
	}
	defer server.Shutdown(10 * time.Second)
	return server.Run(ctx)
}