package getCmd

import (
	"context"
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

var (
	USERSEARCH       = "looking for user"
	FAILEDUSERSEARCH = "failed to get user"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "get a user from the db",
	Long:  `get a user from the db by name, email, username, or role`,
	RunE: func(cmd *cobra.Command, args []string) error {
		all, err := cmd.Flags().GetBool("all")

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return fmt.Errorf("name flag error", "name", name, "error", err)
		}

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			return fmt.Errorf("username flag error", "username", username, "error", err)
		}

		email, err := cmd.Flags().GetString("email")
		if err != nil {
			return fmt.Errorf("email flag error", "email", email, "error", err)
		}

		role, err := cmd.Flags().GetString("role")
		if err != nil {
			return fmt.Errorf("role flag error", "role", role, "error", err)
		}

		uid, err := cmd.Flags().GetString("uid")
		if err != nil {
			return fmt.Errorf("uid flag error", "uid", uid, "error", err)
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
				log.Errorf("failed getting all users", "error", err)
				return nil
			}
		case email != "":
			if err := getUserByEmail(log, ctx, s, email); err != nil {
				log.Errorf(FAILEDUSERSEARCH, "eamil", email, "err", err)
				return nil
			}
		case username != "":
			if err := getUsersByUsername(log, ctx, s, username); err != nil {
				log.Errorf(FAILEDUSERSEARCH, "username", username, "err", err)
				return nil
			}
		case name != "":
			if err := getUserByName(log, ctx, s, name); err != nil {
				log.Errorf(FAILEDUSERSEARCH, "name", name, "err", err)
				return nil
			}
		case role != "":
			if err := getUserByRole(log, ctx, s, role); err != nil {
				log.Errorf(FAILEDUSERSEARCH, "role", role, "error", err)
				return nil
			}
		case uid != "":
			if err := getUserByUID(log, ctx, s, uid); err != nil {
				log.Errorf(FAILEDUSERSEARCH, "uid", uid, "error", err)
				return nil
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
	log.Info("getting all users")
	usrs, err := s.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	if err := displayUsers(usrs); err != nil {
		return err
	}

	return nil
}

func getUserByEmail(log *log.Logger, ctx context.Context, s *user.Store, email string) error {
	log.Info(USERSEARCH, "email", email)
	usrs, err := s.GetUsersByEmail(ctx, email, false)
	if err != nil {
		return err
	}

	if err := displayUsers(usrs); err != nil {
		return err
	}

	return nil
}

func getUsersByUsername(log *log.Logger, ctx context.Context, s *user.Store, uname string) error {
	log.Info(USERSEARCH, "username", uname)
	usrs, err := s.GetUsersByUsername(ctx, uname, false)
	if err != nil {
		return err
	}

	if err := displayUsers(usrs); err != nil {
		return err
	}

	return nil
}

func getUserByName(log *log.Logger, ctx context.Context, s *user.Store, name string) error {
	log.Info(USERSEARCH, "name", name)
	usrs, err := s.GetUsersByName(ctx, name, false)
	if err != nil {
		return err
	}

	if err := displayUsers(usrs); err != nil {
		return err
	}

	return nil
}

func getUserByUID(log *log.Logger, ctx context.Context, s *user.Store, uid string) error {
	log.Info(USERSEARCH, "UID", uid)
	usr, err := s.GetUserByUID(ctx, uid)
	if err != nil {
		return err
	}

	if err := displayUsers([]models.User{usr}); err != nil {
		return err
	}

	return nil
}

func getUserByRole(log *log.Logger, ctx context.Context, s *user.Store, role string) error {
	log.Info(USERSEARCH, "role", role)
	usrs, err := s.GetUsersByRole(ctx, role)
	if err != nil {
		return err
	} else if len(usrs) == 0 {
		log.Infof("no users found", "role", role, "num_users", len(usrs))
		return nil
	}

	if err := displayUsers(usrs); err != nil {
		return err
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
