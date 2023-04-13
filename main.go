package main

import (
	"log"
	"os"

	"github.com/impzero/creeder/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("error: %v", err)
		os.Exit(1)
	}
}
