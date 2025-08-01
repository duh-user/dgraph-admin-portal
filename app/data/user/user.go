// Package user holds the types and functions for
// storing, seraching, and adding new users
package user

import (
	"context"
	"dgraph-client/data/models"
	"dgraph-client/data/role"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"golang.org/x/crypto/bcrypt"
)

// Errors
var (
	ErrNoExists     = errors.New("user does not exit")
	ErrExists       = errors.New("user exists")
	ErrNotFound     = errors.New("user not found")
	ErrPassNotMatch = errors.New("passwords do not match")
)

var query string

// Store will manage the user store API's
type Store struct {
	log *log.Logger
	dgo *dgo.Dgraph
}

// NewStore starts a new db store
func NewStore(log *log.Logger, dgo *dgo.Dgraph) *Store {
	return &Store{
		log: log,
		dgo: dgo,
	}
}

// Add will add a new user to the db if the user doesn't already exist
// if the user existss the found user is returned
// if added the user with uid is returned
func (s *Store) Add(ctx context.Context, newUser *models.NewUser, now time.Time) (models.User, error) {
	nullUsr := models.User{}

	if usrs, err := s.GetUsersByEmail(ctx, newUser.Email, true); err == nil && len(usrs) > 0 {
		for _, usr := range usrs {
			if usr.Email == newUser.Email {
				s.log.Infof("user with email %s already exists (UID: %s)", newUser.Email, usr.UID)
				return usr, ErrExists
			}
		}
	} else if err != nil && !errors.Is(err, ErrNotFound) {
		return nullUsr, fmt.Errorf("failed to check email in db - %w", err)
	}

	if usrs, err := s.GetUsersByUsername(ctx, newUser.UserName, true); err == nil && len(usrs) > 0 {
		for _, usr := range usrs {
			if usr.UserName == newUser.UserName {
				s.log.Infof("user with username %s already exists (UID: %s)", newUser.UserName, usr.UID)
				return usr, ErrExists
			}
		}
	} else if err != nil && !errors.Is(err, ErrNotFound) {
		return nullUsr, fmt.Errorf("failed to check username in db - %w", err)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(newUser.Pass), bcrypt.DefaultCost)
	if err != nil {
		return nullUsr, fmt.Errorf("error hashing pass - %v", err)
	}

	rs := role.NewStore(s.log, s.dgo)
	gotRole, err := rs.GetRoleByName(ctx, newUser.Role)
	if err != nil {
		if errors.Is(err, role.ErrNotFound) {
			return nullUsr, fmt.Errorf("role %s not found %w", newUser.Role, err)
		}
		return nullUsr, fmt.Errorf("error getting role %s - %w", newUser.Role, err)

	}

	user := models.User{
		UserName:     newUser.UserName,
		Name:         newUser.Name,
		Email:        newUser.Email,
		Role:         []models.Role{gotRole},
		PassHash:     string(passHash),
		DateCreated:  now,
		LastSeen:     now,
		LastModified: now,
	}
	r, err := rs.GetRoleByName(ctx, newUser.Role)
	if err != nil {
		return nullUsr, fmt.Errorf("role not found %v", err)
	}

	user.Role = []models.Role{r}

	return s.add(ctx, user)
}

// GetUserByName return user found by provided name
func (s *Store) GetUsersByName(ctx context.Context, name string, exact bool) ([]models.User, error) {
	vars := make(map[string]string)
	vars["$name"] = name

	if exact {
		query = QBYNAMEEXACT
	} else {
		query = QBYNAMEFUZZY
	}

	usrs, err := s.queryUser(ctx, query, vars)
	if err != nil {
		return []models.User{}, err
	}

	return usrs, nil
}

// GetUserByUsername return user found by provided username
func (s *Store) GetUsersByUsername(ctx context.Context, username string, exact bool) ([]models.User, error) {
	vars := make(map[string]string)
	vars["$user_name"] = username

	if exact {
		query = QBYUNAMEEXACT
	} else {
		query = QBYUNAMEFUZZY
	}

	usrs, err := s.queryUser(ctx, query, vars)
	if err != nil {
		return []models.User{}, err
	}

	return usrs, nil
}

// GetUserByEmail returns user found by provided email
func (s *Store) GetUsersByEmail(ctx context.Context, email string, exact bool) ([]models.User, error) {
	vars := make(map[string]string)
	vars["$email"] = email

	if exact {
		query = QBYEMAILEXACT
	} else {
		query = QBYEMAILFUZZY
	}

	usrs, err := s.queryUser(ctx, query, vars)
	for _, usr := range usrs {
		fmt.Println(usr.Name)
	}
	if err != nil {
		return []models.User{}, err
	}

	return usrs, nil
}

// GetUserByUID return user found by proided uid
func (s *Store) GetUserByUID(ctx context.Context, uid string) (models.User, error) {
	vars := make(map[string]string)
	vars["$uid"] = uid
	query = QBYUID

	usr, err := s.queryUser(ctx, query, vars)
	if err == nil && len(usr) < 1 {
		return models.User{}, ErrNotFound
	} else if err != nil {
		return models.User{}, err
	}

	return usr[0], nil
}

// GetUserByRole return all users for a proided role
func (s *Store) GetUsersByRole(ctx context.Context, role string) ([]models.User, error) {
	vars := make(map[string]string)
	vars["$role"] = role
	query = QBYROLE

	roles, err := s.queryUserWithRole(ctx, query, vars)
	if err == nil && len(roles) < 1 {
		return []models.User{}, ErrNotFound
	} else if err != nil {
		return []models.User{}, err
	}

	fmt.Println(roles)

	var usrs []models.User
	for _, role := range roles {
		for _, usr := range role.ReverseEdge {
			fmt.Println(usr.Name)
			usrs = append(usrs, usr)
		}
	}

	return usrs, nil
}

// GetAllUsers returns all users including admins
func (s *Store) GetAllUsers(ctx context.Context) ([]models.User, error) {
	query = QALLUSERS

	usrs, err := s.queryUser(ctx, query, nil)
	if err == nil && len(usrs) < 1 {
		return []models.User{}, ErrNotFound
	} else if err != nil {
		return []models.User{}, err
	}

	return usrs, nil
}

// UpdateUser updates a user in the store
func (s *Store) Update(ctx context.Context, usr models.User) error {
	if usr.UID == "" {
		return fmt.Errorf("missing UID")
	}

	if _, err := s.GetUserByUID(ctx, usr.UID); err != nil {
		return ErrNoExists
	}

	return s.update(ctx, usr)
}

// DeleteUser deletes a user from the store
func (s *Store) Delete(ctx context.Context, usr models.User) error {
	if usr.UID == "" {
		return fmt.Errorf("missing UID")
	}

	if _, err := s.GetUserByUID(ctx, usr.UID); err != nil {
		return ErrNoExists
	}

	return s.delete(ctx, usr.UID)
}

// ------ //

// add uses the client to add a user
func (s *Store) add(ctx context.Context, usr models.User) (models.User, error) {
	jsonUser, err := json.Marshal(usr)
	if err != nil {
		return models.User{}, fmt.Errorf("unable to marshal user to json - %v", err)
	}

	mu := &api.Mutation{
		SetJson:   jsonUser,
		CommitNow: true,
	}

	s.log.Infof("request to add user - %s", usr.UserName)

	resp, err := s.dgo.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return models.User{}, fmt.Errorf("unable to add user to db - %v", err)
	}

	if len(resp.Uids) == 0 {
		return models.User{}, fmt.Errorf("user id not returned - %v", resp.Json)
	}

	s.log.Infof("user added - %s", usr.UID)

	usr.UID = resp.Uids["0"]
	return usr, nil
}

func (s *Store) queryUserWithRole(ctx context.Context, q string, vars map[string]string) ([]models.Role, error) {
	s.log.Infof("request to query user with role - %s", q)
	resp, err := s.dgo.NewTxn().QueryWithVars(ctx, q, vars)
	if err != nil {
		return []models.Role{}, fmt.Errorf("dgo tx failed - QueryWithVars - %v", err)
	}

	fmt.Println(string(resp.Json))

	type Response struct {
		Roles []models.Role `json:"query"`
	}

	var r Response
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return []models.Role{}, fmt.Errorf("error while unmarshaling query result - %v", err)
	}

	fmt.Println(r)

	if len(r.Roles) < 1 {
		return []models.Role{}, ErrNotFound
	}

	return r.Roles, nil
}

func (s *Store) queryUser(ctx context.Context, q string, vars map[string]string) ([]models.User, error) {
	s.log.Infof("request to query user - %s", q)
	resp, err := s.dgo.NewTxn().QueryWithVars(ctx, q, vars)
	if err != nil {
		return []models.User{}, fmt.Errorf("dgo tx failed - QueryWithVars - %v", err)
	}

	type Response struct {
		Users []models.User `json:"query"`
	}

	var r Response
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return []models.User{}, fmt.Errorf("error while unmarshaling query result - %v", err)
	}

	if len(r.Users) < 1 {
		fmt.Println(len(r.Users))
		return []models.User{}, ErrNotFound
	}

	s.log.Infof("returned %d users", len(r.Users))

	return r.Users, nil
}

func (s *Store) update(ctx context.Context, usr models.User) error {
	mutation := &api.Mutation{
		CommitNow: true,
	}

	jsonUser, err := json.Marshal(usr)
	if err != nil {
		return fmt.Errorf("unable to marshal user to json - %v", err)
	}

	s.log.Infof("request to update user - %s", usr.UID)
	mutation.SetJson = jsonUser
	resp, err := s.dgo.NewTxn().Mutate(ctx, mutation)
	if err != nil {
		return fmt.Errorf("error updating user - %v", err)
	}

	if len(resp.Uids) != 1 {
		return fmt.Errorf("failed updating user\nReturned UIDs: %d", len(resp.Uids))
	}

	s.log.Infof("user updated successfully - %s", usr.UID)

	return nil
}

func (s *Store) delete(ctx context.Context, usrID string) error {
	mutation := &api.Mutation{
		CommitNow: true,
	}

	delUser := map[string]string{"uid": usrID}
	jsonUser, err := json.Marshal(delUser)
	if err != nil {
		return fmt.Errorf("failed to marshal empty user - %v", err)
	}

	mutation.DeleteJson = jsonUser

	s.log.Infof("request to delete user : %s", usrID)
	_, err = s.dgo.NewTxn().Mutate(ctx, mutation)
	if err != nil {
		return fmt.Errorf("unable to delete user - %v", err)
	}

	s.log.Infof("%s : %s", "user deleted", usrID)

	return nil
}
