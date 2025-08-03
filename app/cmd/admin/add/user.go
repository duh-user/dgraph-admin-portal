package addCmd

import (
	"context"
	"dgraph-client/data"
	"dgraph-client/data/models"
	"dgraph-client/data/user"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "add a user to the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		log := log.New(os.Stdout)
		traceID := uuid.New().String()
		log.SetPrefix(traceID)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		dgc, cncl := data.NewDGClient(cfg)
		defer cncl()

		s := user.NewStore(log, dgc.Client)

		usr, err := initFlags(cmd)
		if err != nil {
			return fmt.Errorf("unable to init flags - %w", err)
		}

		if err := addUser(log, ctx, s, usr); err != nil {
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

func addUser(log *log.Logger, ctx context.Context, s *user.Store, newUsr *models.NewUser) error {
	usr, err := s.Add(ctx, newUsr, time.Now())
	if err != nil {
		return err
	}

	log.Info("user created successfully - ", usr.UID)

	return nil
}

func initFlags(cmd *cobra.Command) (*models.NewUser, error) {
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return nil, err
	}

	username, err := cmd.Flags().GetString("username")
	if err != nil {
		return nil, err
	}

	email, err := cmd.Flags().GetString("email")
	if err != nil {
		return nil, err
	}

	pass, err := cmd.Flags().GetString("password")
	if err != nil {
		return nil, err
	}

	role, err := cmd.Flags().GetString("role")
	if err != nil {
		return nil, err
	}

	usr := models.NewUser{
		Name:     name,
		UserName: username,
		Email:    email,
		Pass:     pass,
		Role:     role,
	}

	return &usr, nil
}
