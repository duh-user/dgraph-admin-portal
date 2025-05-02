package cmd

import (
	adminCmd "dgraph-client/cmd/admin"
	apiCmd "dgraph-client/cmd/api"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "dgraph-client",
	Short: "testing for dgraph backend",
	Long:  `this is a test to see how I would build a program using dgraph for by backend db`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	initConfig()
	rootCmd.AddCommand(adminCmd.Cmd)
	rootCmd.AddCommand(apiCmd.Cmd)
}

func initConfig() {
	viper.AutomaticEnv()
}
