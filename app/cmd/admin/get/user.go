package getCmd

import (
	"context"
	"dgraph-client/config"
	"dgraph-client/data"
	"dgraph-client/data/user"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var usr user.User

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "get a user from the db",
	Long:  `get a user from the db by name, email, username, or role`,
	Run: func(cmd *cobra.Command, args []string) {
		log := log.New(os.Stdout, "ADMINCMD - ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

		switch {
		case usr.Email != "":
			if err := getUserByEmail(log, cfg); err != nil {
				log.Fatalf("unable get user %s - %v", usr.Email, err)
			}
		case usr.UserName != "":
			if err := getUsersByUsername(log, cfg); err != nil {
				log.Fatalf("unable get user %s - %v", usr.UserName, err)
			}
		case usr.Name != "":
			if err := getUserByName(log, cfg); err != nil {
				log.Fatalf("unable to get user %s - %v", usr.Name, err)
			}
		default:
			fmt.Println("Please enter username, name, email or role")
			cmd.Help()
			os.Exit(1)
		}

	},
}

func init() {
	userCmd.Flags().StringVar(&usr.Name, "name", "", "full name of the user")
	userCmd.Flags().StringVar(&usr.UserName, "username", "", "username of the user")
	userCmd.Flags().StringVar(&usr.Email, "email", "", "email address of the user")
}

func getUserByEmail(log *log.Logger, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dgc, cncl := data.NewDGClient(cfg)
	defer cncl()

	s := user.NewStore(log, dgc.Client)

	traceID := uuid.New().String()
	log.Printf("looking for users matching %s", usr.Email)
	usrs, err := s.GetUsersByEmail(ctx, traceID, usr.Email)
	if err != nil {
		return err
	}

	for _, usr := range usrs {
		log.Println("user found: - ", usr.UID)
		fmt.Printf("UID: %s\nUsername: %s\nName: %s\n", usr.UID, usr.UserName, usr.Name)
		for _, role := range usr.Role {
			fmt.Printf("Role: %s\n", role.Name)
		}
	}

	return nil
}

func getUsersByUsername(log *log.Logger, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dgc, cncl := data.NewDGClient(cfg)
	defer cncl()

	s := user.NewStore(log, dgc.Client)

	traceID := uuid.New().String()
	usrs, err := s.GetUsersByUsername(ctx, traceID, usr.UserName)
	if err != nil {
		return err
	}

	for _, usr := range usrs {
		log.Println("user found - ", usr.UID)
	}

	return nil
}

func getUserByName(log *log.Logger, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dgc, cncl := data.NewDGClient(cfg)
	defer cncl()

	s := user.NewStore(log, dgc.Client)

	traceID := uuid.New().String()
	usrs, err := s.GetUsersByName(ctx, traceID, usr.Name)
	if err != nil {
		return err
	}

	fmt.Println(len(usrs), "users found:")
	for _, usr := range usrs {
		fmt.Printf("UID: %s\nUsername: %s\nName: %s\nRole: %s\n", usr.UID, usr.UserName, usr.Name, usr.Role[0].Name)
	}

	return nil
}
