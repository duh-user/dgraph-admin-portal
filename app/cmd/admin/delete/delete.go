package delete

import (
	"dgraph-client/config"

	"github.com/spf13/cobra"
)

var cfg *config.Config

var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "delete shit",
	Long:  `delete the schema, db, users, etc`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	cfg = config.InitConfig()
	//Cmd.AddCommand(dataCmd)
	Cmd.AddCommand(everythingCmd)
}
