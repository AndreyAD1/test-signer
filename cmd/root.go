package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "test-signer",
	Short: "The 'Test Signer' service.",
	Long: `The Test signer is a service that accepts a set of answers and 
questions and signs that the user has finished the " test " at this point in time. 
The signatures are stored and can later be verified by a different service.`,
}

func Execute() error {
	return RootCmd.Execute()
}