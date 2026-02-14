package main

import (
	"os"

	"github.com/Go-Yadro-Group-1/Jira-Connector/cmd/internal/cli"
)

func main() {
	rootCmd := cli.NewRootCmd()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
