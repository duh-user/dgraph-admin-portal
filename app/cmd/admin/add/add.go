package addCmd

import (
	"dgraph-client/config"

	"github.com/spf13/cobra"
)

var cfg *config.Config

var Cmd = &cobra.Command{
	Use:   "add",
	Short: "add a user or data to the database",
}

func init() {
	cfg = config.InitConfig()
	Cmd.AddCommand(userCmd)
}
