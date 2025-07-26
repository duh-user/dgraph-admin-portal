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

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "get a user from the db",
	Long:  `get a user from the db by name, email, username, or role`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log := log.New(os.Stdout, "ADMINCMD - ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return fmt.Errorf("unable to get name flag - %v", err)
		}

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return fmt.Errorf("unable to get username flag - %v", err)
		}

		email, err := cmd.Flags().GetString("email")
		if err != nil {
			return fmt.Errorf("unable to get email flag - %v", err)
		}

		role, err := cmd.Flags().GetString("role")
		if err != nil {
			return fmt.Errorf("unable to get role flag - %v", err)
		}

		uid, err := cmd.Flags().GetString("uid")
		if err != nil {
			return fmt.Errorf("unable to get uid flag - %v", err)
		}

		switch {
		case email != "":
			if err := getUserByEmail(log, cfg, email); err != nil {
				log.Fatalf("unable get user %s - %v", email, err)
			}
		case username != "":
			if err := getUsersByUsername(log, cfg, username); err != nil {
				log.Fatalf("unable get user %s - %v", username, err)
			}
		case name != "":
			if err := getUserByName(log, cfg, name); err != nil {
				log.Fatalf("unable to get user %s - %w", name, err)
			}
		case role != "":
			if err := getUserByRole(log, cfg, role); err != nil {
				log.Fatalf("unable to get users with role %s - %w", role, err)
			}
		case uid != "":
			if err := getUserByUID(log, cfg, uid); err != nil {
				log.Fatalf("unable to get user with uid %s - %w", uid, err)
			}
		default:
			fmt.Println("error - please use flag --username, --uid, --name, --email or --role")
			cmd.Help()
			return fmt.Errorf("no search criteria provided")
		}
		return nil
	},
}

func init() {
	userCmd.Flags().String("name", "", "full name of the user")
	userCmd.Flags().String("username", "", "username of the user")
	userCmd.Flags().String("email", "", "email address of the user")
	userCmd.Flags().Bool("all", false, "get all users")
	userCmd.Flags().String("role", "", "get all users of a specific role")
	userCmd.Flags().String("uid", "", "get user by uid")
}

/*
func query(log *log.Logger, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dgc, cncl := data.NewDGClient(cfg)
	defer cncl()

	s := user.NewStore(log, dgc.Client)

	traceID := uuid.New().String()
}
*/

func getUserByEmail(log *log.Logger, cfg *config.Config, email string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dgc, cncl := data.NewDGClient(cfg)
	defer cncl()

	s := user.NewStore(log, dgc.Client)

	traceID := uuid.New().String()
	log.Printf("looking for users matching %s", email)
	usrs, err := s.GetUsersByEmail(ctx, traceID, email, false)
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

func getUsersByUsername(log *log.Logger, cfg *config.Config, uname string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dgc, cncl := data.NewDGClient(cfg)
	defer cncl()

	s := user.NewStore(log, dgc.Client)

	traceID := uuid.New().String()
	usrs, err := s.GetUsersByUsername(ctx, traceID, uname, false)
	if err != nil {
		return err
	}

	for _, usr := range usrs {
		log.Println("user found - ", usr.UID)
	}

	return nil
}

func getUserByName(log *log.Logger, cfg *config.Config, name string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dgc, cncl := data.NewDGClient(cfg)
	defer cncl()

	s := user.NewStore(log, dgc.Client)

	traceID := uuid.New().String()
	usrs, err := s.GetUsersByName(ctx, traceID, name, false)
	if err != nil {
		return err
	}

	fmt.Println(len(usrs), "users found:")
	for _, usr := range usrs {
		fmt.Printf("UID: %s\nUsername: %s\nName: %s\nRole: %s\n", usr.UID, usr.UserName, usr.Name, usr.Role[0].Name)
	}

	return nil
}

func getUserByUID(log *log.Logger, cfg *config.Config, uid string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dgc, cncl := data.NewDGClient(cfg)
	defer cncl()

	s := user.NewStore(log, dgc.Client)

	traceID := uuid.New().String()

	usr, err := s.GetUserByUID(ctx, traceID, uid)
	if err != nil {
		return err
	}

	fmt.Println("Found user:")
	fmt.Printf("UID: %s\nUsername: %s\nName: %s\nRole: %s\n", usr.UID, usr.UserName, usr.Name, usr.Role[0].Name)

	return nil
}

func getUserByRole(log *log.Logger, cfg *config.Config, role string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dgc, cncl := data.NewDGClient(cfg)
	defer cncl()

	s := user.NewStore(log, dgc.Client)

	traceID := uuid.New().String()

	usrs, err := s.GetUsersByRole(ctx, traceID, role)
	if err != nil {
		return err
	}

	fmt.Println(len(usrs), " users found:")
	for _, usr := range usrs {
		fmt.Printf("UID: %s\nUsername: %s\nName: %s\nRole: %s\n", usr.UID, usr.UserName, usr.Name, usr.Role[0].Name)
	}

	return nil
}

func getAllUsers(log *log.Logger, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dgc, cncl := data.NewDGClient(cfg)
	defer cncl()

	s := user.NewStore(log, dgc.Client)

	traceID := uuid.New().String()

	usrs, err := s.GetAllUsers(ctx, traceID)
	if err != nil {
		return err
	}

	fmt.Println(len(usrs), "users found:")
	for _, usr := range usrs {
		fmt.Printf("UID: %s\nUsername: %s\nName: %s\nRole: %s\n", usr.UID, usr.UserName, usr.Name, usr.Role[0].Name)
	}

	return nil
}
