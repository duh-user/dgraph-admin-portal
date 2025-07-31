package addCmd

import (
	"context"
	"dgraph-client/config"
	"dgraph-client/data"
	"dgraph-client/data/models"
	"dgraph-client/data/user"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "add a user to the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		log := log.New(os.Stdout, "ADMINCMD - ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

		usr, err := initFlags(cmd)
		if err != nil {
			return fmt.Errorf("unable to init flags - %w", err)
		}

		if err := addUser(log, cfg, usr); err != nil {
			return fmt.Errorf("unable to add user - %w", err)
		}
		return nil
	},
}

func init() {
	userCmd.Flags().String("name", "", "full name of the new user")
	userCmd.Flags().String("username", "", "full name of the new user")
	userCmd.Flags().String("email", "", "email address of the new user")
	userCmd.Flags().String("password", "", "password of the new user")
	userCmd.Flags().String("role", "user", "user role - default: user")
	userCmd.MarkFlagRequired("name")
	userCmd.MarkFlagRequired("email")
	userCmd.MarkFlagRequired("password")
	userCmd.MarkFlagRequired("username")
	userCmd.MarkFlagRequired("role")
}

func addUser(log *log.Logger, cfg *config.Config, newUsr *models.NewUser) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dgc, cncl := data.NewDGClient(cfg)
	defer cncl()

	s := user.NewStore(log, dgc.Client)

	traceID := uuid.New().String()
	usr, err := s.Add(ctx, traceID, newUsr, time.Now())
	if err != nil {
		return err
	}

	log.Println("user created successfully - ", usr.UID)

	return nil
}

func initFlags(cmd *cobra.Command) (*models.NewUser, error) {
	usr := models.NewUser{}

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return &usr, err
	}

	username, err := cmd.Flags().GetString("username")
	if err != nil {
		return &usr, err
	}

	email, err := cmd.Flags().GetString("email")
	if err != nil {
		return &usr, err
	}

	pass, err := cmd.Flags().GetString("password")
	if err != nil {
		return &usr, err
	}

	role, err := cmd.Flags().GetString("role")
	if err != nil {
		return &usr, err
	}

	usr = models.NewUser{
		Name:     name,
		UserName: username,
		Email:    email,
		Pass:     pass,
		Role:     role,
	}

	return &usr, nil
}
