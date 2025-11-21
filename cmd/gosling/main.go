package main

import (
	"fmt"
	"os"

	"github.com/ElishaFlacon/fast-sobes-auth/pkg/gosling"
)

func main() {
	if err := gosling.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
