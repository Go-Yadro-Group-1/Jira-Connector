/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"fmt"

	"github.com/Go-Yadro-Group-1/Jira-Connector/cmd/config"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the Jira Connector server",
	Long:  "Starts the Jira Connector server with configured services",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()

		fmt.Println("Starting Jira Connector")
		fmt.Printf("Jira URL: %s\n", cfg.Jira.URL)
		fmt.Printf("Database: %s@%s:%d/%s\n", cfg.DB.User, cfg.DB.Host, cfg.DB.Port, cfg.DB.DBName)

		if cfg.Broker.URL != "" {
			fmt.Printf("Broker: %s\n", cfg.Broker.URL)
		}

		fmt.Println("Service started")

		select {}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
