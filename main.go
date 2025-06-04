package main

import (
	"fmt"
	"os"

	"github.com/un-versed/go-sbv-to-srt/cmd"
)

var (
	Version = "dev"
)

func main() {
	cmd.SetVersionInfo(Version)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
