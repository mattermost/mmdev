package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/mattermost/mmdev/internal/node"
	"github.com/mattermost/mmdev/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setupCmd = &cobra.Command{
	Use:          "setup",
	Short:        "Setup development environment",
	SilenceUsage: true,
}

var mobileCmd = &cobra.Command{
	Use:          "mobile",
	Short:        "Setup the environment to develop the mobile app",
	RunE:         mobileCmdF,
	SilenceUsage: true,
}

func mobileCmdF(c *cobra.Command, args []string) error {
	cfg, err := resolveConfig(c.OutOrStdout())
	if err != nil {
		return err
	}

	var result *multierror.Error

	// for _, repo := range cfg.Repositories {
	// skip if the repository is not there
	// fmt.Println("scanning repository %s...", repo)
	// if _, err := os.Stat(repo.Path); os.IsNotExist(err) {
	// 	continue
	// }

	// remote := repo.Remote
	// if len(args) > 0 {
	// 	remote = args[0]
	// }

	// cmdArgs := append([]string{"pull"}, args...)

	// err := git.RunGitCommand(cfg, repo, cmdArgs...)
	// if err != nil {
	// 	result = multierror.Append(result, err)
	// }

	// fmt.Fprintf(c.OutOrStdout(), "%s: successfully pulled from %s\n", repo.Name, remote)
	// }

	nvmError := node.InstallNVMIfNeeded(cfg.NVM)
	if nvmError != nil {
		result = multierror.Append(result, nvmError)
	}

	if nvmError == nil {
		nodeError := node.InstallNodeIfNeeded(cfg.NodeJS)
		if nodeError != nil {
			result = multierror.Append(result, nodeError)
		}
	}

	return result.ErrorOrNil()

}

func resolveConfig(w io.Writer) (*model.Config, error) {
	if !viper.IsSet("config") {
		fmt.Fprintln(w, "no config detected, continuing with defaults...")
		return model.DefaultConfig(), nil
	}

	fileContents, err := os.ReadFile(viper.GetString("config"))
	if err != nil {
		return nil, fmt.Errorf("there was a problem reading the config file: %w", err)
	}

	var cfg model.Config
	if err := json.Unmarshal(fileContents, &cfg); err != nil {
		return nil, fmt.Errorf("there was a problem parsing the config file: %w", err)
	}

	return &cfg, nil
}
