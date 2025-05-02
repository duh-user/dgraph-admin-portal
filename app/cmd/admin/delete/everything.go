package delete

import (
	"context"
	"dgraph-client/config"
	"dgraph-client/data"
	"dgraph-client/data/schema"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var everythingCmd = &cobra.Command{
	Use:   "everything",
	Short: "blow it all the fuck up",
	Long:  `delete the schema and all data...start anew`,
	Run: func(cmd *cobra.Command, args []string) {
		log := log.New(os.Stdout, "ADMINCMD - ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

		if err := deleteEverything(log, cfg); err != nil {
			log.Fatal("error while killing everything - ", err)
		}

	},
}

func deleteEverything(log *log.Logger, cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dgc, cncl := data.NewDGClient(cfg)
	defer cncl()

	/*
		if err := dgc.HealthCheck(ctx, 5*time.Second); err != nil {
			return fmt.Errorf("waiting for db... - ", err)
		}
	*/

	schema, err := schema.NewSchema(dgc.Client)
	if err != nil {
		return fmt.Errorf("preping schema... - %v", err)
	}

	if err := schema.DropAll(ctx); err != nil {
		return fmt.Errorf("killing everything... - %v", err)
	}

	return nil
}
