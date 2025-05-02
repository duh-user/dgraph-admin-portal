package apiCmd

import (
	"dgraph-client/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfg    *config.Config
	apiCfg *config.APIConfig
)

var Cmd = &cobra.Command{
	Use:   "api",
	Short: "control the api server",
	Long:  `start and stop the api server`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg = config.InitConfig()
		apiCfg = cfg.InitAPIConfig()
	},
}

func init() {
	Cmd.PersistentFlags().String("api-addr", "localhost:8888",
		"set hostname for service. default: localhost:8081")
	viper.BindPFlag("api-addr", Cmd.PersistentFlags().Lookup("api-addr"))
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
