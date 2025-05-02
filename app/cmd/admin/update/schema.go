package update

import (
	"context"
	"dgraph-client/config"
	"dgraph-client/data"
	"dgraph-client/data/schema"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "update the schema",
	Long:  `update schema from hardcoded....tbd from backup, file, etc`,
	Run: func(cmd *cobra.Command, args []string) {
		log := log.New(os.Stdout, "ADMINCMD - ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

		if err := updateSchema(log, cfg); err != nil {
			log.Println("error during schema creation - ", err)
		}
	},
}

/*
var init() {
	// TODO add file declaration
	// initSchema.PersistentFlags().StringVar(&newUser.Name, "file", "", "schema file location")
}
*/

func updateSchema(log *log.Logger, cfg *config.Config) error {
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
		return fmt.Errorf("error preping schema... - %v", err)
	}

	traceID := uuid.New().String()

	if err := schema.InitSchema(ctx); err != nil {
		return fmt.Errorf("error creating schema... - %v", err)
	}

	log.Println("schema updated successfully")

	if err := schema.InitRoles(ctx, log, traceID); err != nil {
		return fmt.Errorf("error creating roles... - %v", err)
	}

	return nil
}
