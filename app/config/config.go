package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// TODO - add TLS support between containers and REST API
	// CertsDir   string
	DGAddr string
}

type APIConfig struct {
	ApiAddr         string
	ApiReadTimeout  time.Duration
	ApiWriteTimeout time.Duration
	ApiIdleTimeout  time.Duration
	DGAddr          string
	// TODO - add TLS support
}

func InitConfig() *Config {
	if err := pullConfig(); err != nil {
		log.Fatalf("fatal error reading config file -", err)
	}
	cfg := &Config{
		DGAddr: viper.GetString("DGADDR"),
	}

	return cfg

}

func (c *Config) InitAPIConfig() *APIConfig {
	apiCfg := &APIConfig{
		ApiAddr:         viper.GetString("APIADDR"),
		ApiReadTimeout:  time.Second * 5,
		ApiWriteTimeout: time.Second * 10,
		ApiIdleTimeout:  time.Second * 120,
		DGAddr:          c.DGAddr,
	}

	return apiCfg

	// TODO - add TLS support between containers and REST API
	// c.TLSCert = fmt.Sprintf("%s/app.crt", c.CertsDir)
	// c.TLSKey = fmt.Sprintf("%s/app.key", c.CertsDir)
}

func pullConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/dgraph-client/")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("config file not found. using defaults -", err)
			//			loadDefaults()
		} else {
			return err
		}
	}
	viper.AutomaticEnv()

	return nil
}
