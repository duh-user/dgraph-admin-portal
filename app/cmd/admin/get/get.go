package getCmd

import (
	"dgraph-client/config"

	"github.com/spf13/cobra"
)

var cfg *config.Config

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "get data from the db",
}

func init() {
	cfg = config.InitConfig()
	Cmd.AddCommand(userCmd)
}
