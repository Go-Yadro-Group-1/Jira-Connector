package cli

import (
	"context"
	_ "embed"
	"errors"
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
	flagConfigFile    = "config"
	flagConfigFileH   = "c"
	defaultConfigPath = "config/release.yaml"
)

var (
	errGetConfigFlag   = errors.New("get config path")
	errReadConfig      = errors.New("read config file")
	errValidateConfig  = errors.New("validate config")
	errInitApplication = errors.New("create new fs manager")
)

func NewRunCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"start", "launch"},
		Short:   "Run File System Manager",
		Long:    embeddedRunLongData,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
			defer cancel()

			cfg, err := config.New()
			if err != nil {
				return errors.Join(errReadConfig, err)
			}

			connector, err := app.New(cfg)
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
}
