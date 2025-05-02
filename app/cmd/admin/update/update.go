package update

import (
	"dgraph-client/config"

	"github.com/spf13/cobra"
)

var cfg *config.Config

var Cmd = &cobra.Command{
	Use:   "update",
	Short: "update shit",
	Long:  `update the schema, db, etc....from hardcoded now. TODO from backup/file/etc`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	cfg = config.InitConfig()
	Cmd.AddCommand(schemaCmd)
}
