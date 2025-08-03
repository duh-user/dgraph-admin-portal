package adminCmd

import (
	addCmd "dgraph-client/cmd/admin/add"
	deleteCmd "dgraph-client/cmd/admin/delete"
	getCmd "dgraph-client/cmd/admin/get"
	updateCmd "dgraph-client/cmd/admin/update"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &cobra.Command{
	Use:   "admin",
	Short: "administration commands",
	Long:  `add data, remove data, query data, be an admin`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	Cmd.AddCommand(addCmd.Cmd)
	Cmd.AddCommand(deleteCmd.Cmd)
	Cmd.AddCommand(updateCmd.Cmd)
	Cmd.AddCommand(getCmd.Cmd)
	Cmd.PersistentFlags().String("dg-addr", "localhost:9080",
		"set dgraph host url. default: localhost:9080")
	viper.BindPFlag("dg-addr", Cmd.PersistentFlags().Lookup("dg-addr"))
	// TODO - add TLS support
	//apiCmd.PersistentFlags().String("dg-addr", "https://localhost:9080",
	//	"set dgraph host url. default: localhost:9080")
	//apiCmd.PersistentFlags().String("api-addr", "https://localhost:8081",
	//	"set hostname for service. default: localhost:8081")
	// flag.StringVar(&cfg.CertsDir, "certs-dir", "./certs", "set directory for TLS certs. default ./certs")

}
