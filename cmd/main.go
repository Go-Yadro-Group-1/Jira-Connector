package main

import (
	"os"

	"github.com/Go-Yadro-Group-1/Jira-Connector/cmd/internal/cli"
)

func main() {
	rootCmd := cli.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
