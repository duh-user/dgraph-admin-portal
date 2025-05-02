package addCmd

import (
	"context"
	"dgraph-client/config"
	"dgraph-client/data"
	"dgraph-client/data/user"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "add a user to the database",
	Run: func(cmd *cobra.Command, args []string) {
		log := log.New(os.Stdout, "ADMINCMD - ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

		usr, err := initFlags(cmd)
		if err != nil {
			log.Fatal("unable to init flags -", err)
		}

		if err := addUser(log, cfg, usr); err != nil {
			log.Fatal("unable to add user -", err)
		}
	},
}

func init() {
	userCmd.PersistentFlags().String("name", "", "full name of the new user")
	userCmd.PersistentFlags().String("username", "", "full name of the new user")
	userCmd.PersistentFlags().String("email", "", "email address of the new user")
	userCmd.PersistentFlags().String("password", "", "password of the new user")
	userCmd.PersistentFlags().String("role", "user", "user role - default: user")
	userCmd.MarkPersistentFlagRequired("name")
	userCmd.MarkPersistentFlagRequired("email")
	userCmd.MarkPersistentFlagRequired("password")
	userCmd.MarkPersistentFlagRequired("username")
	userCmd.MarkPersistentFlagRequired("role")
}

func addUser(log *log.Logger, cfg *config.Config, newUsr *user.NewUser) error {
	if newUsr.Name == "" || newUsr.Email == "" || newUsr.Pass == "" || newUsr.Role == "" || newUsr.UserName == "" {
		return Cmd.Usage()
	}

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

func initFlags(cmd *cobra.Command) (*user.NewUser, error) {
	usr := user.NewUser{}

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

	usr = user.NewUser{
		Name:     name,
		UserName: username,
		Email:    email,
		Pass:     pass,
		Role:     role,
	}

	return &usr, nil
}
