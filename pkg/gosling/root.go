package gosling

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gosling",
	Short: "Code generation utility for microservices",
	Long:  `Gosling is a CLI tool for generating and managing microservice components (handlers, usecases, repositories)`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(makeCmd)
	rootCmd.AddCommand(removeCmd)
}
