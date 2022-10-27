package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "mmdev",
	Short: "Utility tool to setup Mattermost developer environments",
}

func init() {
	RootCmd.PersistentFlags().String("config", filepath.Join("$XDG_CONFIG_HOME", "mmdev", "config"), "path to the configuration file")
	_ = viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))

	RootCmd.AddCommand(setupCmd)
	setupCmd.AddCommand(mobileCmd)
}

func main() {
	RootCmd.SetArgs(os.Args[1:])

	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
