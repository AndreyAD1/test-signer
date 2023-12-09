package main

import (
	"os"

	"github.com/AndreyAD1/test-signer/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
