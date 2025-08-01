package getCmd

import (
	"context"
	"dgraph-client/config"
	"dgraph-client/data"
	"dgraph-client/data/models"
	"dgraph-client/data/user"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "get a user from the db",
	Long:  `get a user from the db by name, email, username, or role`,
	RunE: func(cmd *cobra.Command, args []string) error {
		all, err := cmd.Flags().GetBool("all")

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

		log := log.New(os.Stdout)
		traceID := uuid.New().String()
		log.SetPrefix(traceID)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		dgc, cncl := data.NewDGClient(cfg)
		defer cncl()

		s := user.NewStore(log, dgc.Client)

		switch {
		case all:
			if err := getAllUsers(log, ctx, s); err != nil {
				return fmt.Errorf("unable to get all users - %v", err)
			}
		case email != "":
			if err := getUserByEmail(log, cfg, traceID, email); err != nil {
				return fmt.Errorf("unable get user %s - %v", email, err)
			}
		case username != "":
			if err := getUsersByUsername(log, cfg, traceID, username); err != nil {
				return fmt.Errorf("unable get user %s - %v", username, err)
			}
		case name != "":
			if err := getUserByName(log, cfg, traceID, name); err != nil {
				return fmt.Errorf("unable to get user %s - %w", name, err)
			}
		case role != "":
			if err := getUserByRole(log, cfg, traceID, role); err != nil {
				return fmt.Errorf("unable to get users with role %s - %w", role, err)
			}
		case uid != "":
			if err := getUserByUID(log, cfg, traceID, uid); err != nil {
				return fmt.Errorf("unable to get user with uid %s - %w", uid, err)
			}
		default:
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

func getAllUsers(log *log.Logger, ctx context.Context, s *user.Store) error {
	usrs, err := s.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	if err := displayUsers(usrs); err != nil {
		log.Errorf("failed to display users - $v", err)
		return err
	}

	return nil
}

func getUserByEmail(log *log.Logger, cfg *config.Config, traceID string, email string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dgc, cncl := data.NewDGClient(cfg)
	defer cncl()

	s := user.NewStore(log, dgc.Client)

	log.Info("looking for users matching %s", email)
	usrs, err := s.GetUsersByEmail(ctx, email, false)
	if err != nil {
		return err
	}

	if err := displayUsers(usrs); err != nil {
		return err
	}

	return nil
}

func getUsersByUsername(log *log.Logger, cfg *config.Config, traceID string, uname string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dgc, cncl := data.NewDGClient(cfg)
	defer cncl()

	s := user.NewStore(log, dgc.Client)

	usrs, err := s.GetUsersByUsername(ctx, uname, false)
	if err != nil {
		return err
	}

	if err := displayUsers(usrs); err != nil {
		return err
	}

	return nil
}

func getUserByName(log *log.Logger, cfg *config.Config, traceID string, name string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dgc, cncl := data.NewDGClient(cfg)
	defer cncl()

	s := user.NewStore(log, dgc.Client)

	usrs, err := s.GetUsersByName(ctx, name, false)
	if err != nil {
		return err
	}

	if err := displayUsers(usrs); err != nil {
		return err
	}

	return nil
}

func getUserByUID(log *log.Logger, cfg *config.Config, traceID string, uid string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dgc, cncl := data.NewDGClient(cfg)
	defer cncl()

	s := user.NewStore(log, dgc.Client)

	usr, err := s.GetUserByUID(ctx, uid)
	if err != nil {
		return err
	}

	if err := displayUsers([]models.User{usr}); err != nil {
		log.Errorf("unable to display users - %v", err)
		return nil
	}

	return nil
}

func getUserByRole(log *log.Logger, cfg *config.Config, traceID, role string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dgc, cncl := data.NewDGClient(cfg)
	defer cncl()

	s := user.NewStore(log, dgc.Client)

	usrs, err := s.GetUsersByRole(ctx, role)
	if err != nil {
		return err
	} else if len(usrs) == 0 {
		return fmt.Errorf("no users found for role %s", role)
	}

	fmt.Println(len(usrs), " users found:")
	for _, usr := range usrs {
		fmt.Println(usr)
		fmt.Printf("\n-----\nUID: %s\nUsername: %s\nName: %s\nRole: %s\n\n", usr.UID, usr.UserName, usr.Name, usr.Role[0].Name)
	}

	return nil
}

func displayUsers(usrs []models.User) error {
	if len(usrs) < 1 {
		return fmt.Errorf("no users provided to display")
	}

	rows := [][]string{}

	for _, usr := range usrs {
		rows = append(rows, []string{usr.UID, usr.Name, usr.UserName, usr.Email, usr.LastSeen.String(), usr.LastModified.String()})
	}

	var (
		purple    = lipgloss.Color("99")
		gray      = lipgloss.Color("245")
		lightGray = lipgloss.Color("241")

		headerStyle  = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
		cellStyle    = lipgloss.NewStyle().Padding(0, 1).Width(14)
		oddRowStyle  = cellStyle.Foreground(gray)
		evenRowStyle = cellStyle.Foreground(lightGray)
	)

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return headerStyle
			case row%2 == 0:
				return evenRowStyle
			default:
				return oddRowStyle
			}
		}).
		Headers("UID", "name", "username", "email", "last_seen", "last_modified").
		Rows(rows...)

	fmt.Println(t)

	return nil
}
