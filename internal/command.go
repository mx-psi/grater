package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mx-psi/grater/internal/config"
	"github.com/mx-psi/grater/internal/scraper"
)

type Command struct {
	rootCmd    *cobra.Command
	configPath string
}

func (cmd *Command) Execute() error {
	return cmd.rootCmd.Execute()
}

func getDependents(cfg config.Config) ([]scraper.ModuleDep, error) {
	client, err := scraper.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to build client: %w", err)
	}

	// Exactly 1 module is supported.
	return client.Dependents(context.Background(), cfg.Modules[0])
}

func NewCommand() *Command {
	cmd := &Command{}
	cmd.rootCmd = &cobra.Command{
		Use:   "grater",
		Short: "Grater tests changes on Go modules on its dependents",
		RunE: func(*cobra.Command, []string) error {
			cfg, err := config.LoadAndValidate(cmd.configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			deps, err := getDependents(cfg)
			if err != nil {
				return err
			}

			bytes, err := json.MarshalIndent(deps, "", "  ")
			if err != nil {
				return err
			}

			fmt.Printf("%s\n", string(bytes))
			return nil
		},
		SilenceUsage: true,
	}

	cmd.rootCmd.Flags().StringVar(&cmd.configPath, "config", "", "configuration source")
	cmd.rootCmd.MarkFlagRequired("config")
	return cmd
}
