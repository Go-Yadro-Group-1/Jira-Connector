package cli

import (
	"context"
	_ "embed"
	"errors"
<<<<<<< HEAD
	"fmt"
=======
>>>>>>> main
	"log"
	"os"
	"os/signal"

	"github.com/Go-Yadro-Group-1/Jira-Connector/cmd/internal/app"
	"github.com/Go-Yadro-Group-1/Jira-Connector/cmd/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

//go:embed data/run-long-data.txt
var embeddedRunLongData string

const (
<<<<<<< HEAD
	flagProjectKey    = "project"
=======
>>>>>>> main
	flagConfigFile    = "config"
	flagConfigFileH   = "c"
	defaultConfigPath = "config/release.yaml"
)

var (
	errGetConfigFlag   = errors.New("get config path") //nolint:unused
	errReadConfig      = errors.New("read config file")
	errValidateConfig  = errors.New("validate config") //nolint:unused
	errInitApplication = errors.New("create new fs manager")
<<<<<<< HEAD
	errNoProjectKey    = errors.New("flag --project is required")
=======
>>>>>>> main
)

func NewRunCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"start", "launch"},
		Short:   "Run File System Manager",
		Long:    embeddedRunLongData,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error { //nolint:revive
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
			defer cancel()

<<<<<<< HEAD
			cfg, err := config.LoadDevConfig()
=======
			cfg, err := config.New()
>>>>>>> main
			if err != nil {
				return errors.Join(errReadConfig, err)
			}

<<<<<<< HEAD
			projectKey, err := cmd.Flags().GetString(flagProjectKey)
			if err != nil {
				return fmt.Errorf("failed to get project flag: %w", err)
			}
			if projectKey == "" {
				return errNoProjectKey
			}

			connector, err := app.New(cfg.Jira, projectKey)
=======
			connector, err := app.New(cfg)
>>>>>>> main
			if err != nil {
				return errors.Join(errInitApplication, err)
			}

			errChan := connector.Run()
			defer connector.Close()

			select {
			case <-ctx.Done():
				log.Println("stop service by ctx")
			case err := <-errChan:
				return err
			}

			return nil
		},
	}

	initConfigPath(runCmd.Flags())

	return runCmd
}

func initConfigPath(flagset *pflag.FlagSet) {
	flagset.StringP(flagConfigFile, flagConfigFileH, defaultConfigPath, "path to config")
<<<<<<< HEAD
	flagset.String(flagProjectKey, "", "Jira project key to sync")
=======
>>>>>>> main
}
