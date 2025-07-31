// Package user holds the types and functions for
// storing, seraching, and adding new users
package user

import (
	"context"
	"dgraph-client/data/role"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

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
func (s *Store) Add(ctx context.Context, traceID string, newUser *NewUser, now time.Time) (User, error) {
	if usrs, err := s.GetUsersByEmail(ctx, traceID, newUser.Email, true); err == nil && len(usrs) > 0 {
		for _, usr := range usrs {
			if usr.Email == newUser.Email {
				s.log.Printf("%s - user with email %s already exists (UID: %s)", traceID, newUser.Email, usr.UID)
				return usr, ErrExists
			}
		}
	} else if err != nil && !errors.Is(err, ErrNotFound) {
		return User{}, fmt.Errorf("%s - failed to check email in db - %w", traceID, err)
	}

	if usrs, err := s.GetUsersByUsername(ctx, traceID, newUser.UserName, true); err == nil && len(usrs) > 0 {
		for _, usr := range usrs {
			if usr.UserName == newUser.UserName {
				s.log.Printf("%s - user with username %s already exists (UID: %s)", traceID, newUser.UserName, usr.UID)
				return usr, ErrExists
			}
		}
	} else if err != nil && !errors.Is(err, ErrNotFound) {
		return User{}, fmt.Errorf("%s - failed to check username in db - %w", traceID, err)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(newUser.Pass), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("%s : error hashing pass - %v", traceID, err)
	}

	rs := role.NewStore(s.log, s.dgo)
	gotRole, err := rs.GetRoleByName(ctx, traceID, newUser.Role)
	if err != nil {
		if errors.Is(err, role.ErrNotFound) {
			return User{}, fmt.Errorf("%s : role %s not found %w", traceID, newUser.Role, err)
		}
		return User{}, fmt.Errorf("%s : error getting role %s - %w", traceID, newUser.Role, err)

	}

	user := User{
		UserName:     newUser.UserName,
		Name:         newUser.Name,
		Email:        newUser.Email,
		Role:         []role.Role{gotRole},
		PassHash:     string(passHash),
		DateCreated:  now,
		LastSeen:     now,
		LastModified: now,
	}
	r, err := rs.GetRoleByName(ctx, traceID, newUser.Role)
	if err != nil {
		return User{}, fmt.Errorf("%s : role not found %v", traceID, err)
	}

	user.Role = []role.Role{r}

	return s.add(ctx, traceID, user)
}

// GetUserByName return user found by provided name
func (s *Store) GetUsersByName(ctx context.Context, traceID string, name string, exact bool) ([]User, error) {
	vars := make(map[string]string)
	vars["$name"] = name

	if exact {
		query = QBYNAMEEXACT
	} else {
		query = QBYNAMEFUZZY
	}

	usrs, err := s.query(ctx, traceID, query, vars)
	if err != nil {
		return []User{}, err
	}

	return usrs, nil
}

// GetUserByUsername return user found by provided username
func (s *Store) GetUsersByUsername(ctx context.Context, traceID string, username string, exact bool) ([]User, error) {
	vars := make(map[string]string)
	vars["$user_name"] = username

	if exact {
		query = QBYUNAMEEXACT
	} else {
		query = QBYUNAMEFUZZY
	}

	usrs, err := s.query(ctx, traceID, query, vars)
	if err != nil {
		return []User{}, err
	}

	return usrs, nil
}

// GetUserByEmail returns user found by provided email
func (s *Store) GetUsersByEmail(ctx context.Context, traceID string, email string, exact bool) ([]User, error) {
	vars := make(map[string]string)
	vars["$email"] = email

	if exact {
		query = QBYEMAILEXACT
	} else {
		query = QBYEMAILFUZZY
	}

	usrs, err := s.query(ctx, traceID, query, vars)
	for _, usr := range usrs {
		fmt.Println(usr.Name)
	}
	if err != nil {
		return []User{}, err
	}

	return usrs, nil
}

// GetUserByUID return user found by proided uid
func (s *Store) GetUserByUID(ctx context.Context, traceID string, uid string) (User, error) {
	vars := make(map[string]string)
	vars["$uid"] = uid
	query = QBYUID

	usr, err := s.query(ctx, traceID, query, vars)
	if err == nil && len(usr) < 1 {
		return User{}, ErrNotFound
	} else if err != nil {
		return User{}, err
	}

	return usr[0], nil
}

// GetUserByRole return all users for a proided role
func (s *Store) GetUsersByRole(ctx context.Context, traceID string, role string) ([]User, error) {
	vars := make(map[string]string)
	vars["$role"] = role
	query = QBYROLE

	usrs, err := s.query(ctx, traceID, query, vars)
	if err == nil && len(usrs) < 1 {
		return []User{}, ErrNotFound
	} else if err != nil {
		return []User{}, err
	}

	return usrs, nil
}

// GetAllUsers returns all users including admins
func (s *Store) GetAllUsers(ctx context.Context, traceID string) ([]User, error) {
	vars := make(map[string]string)
	vars["$role"] = "admin|user"
	query = QBYROLE

	usrs, err := s.query(ctx, traceID, query, nil)
	if err == nil && len(usrs) < 1 {
		return []User{}, ErrNotFound
	} else if err != nil {
		return []User{}, err
	}

	return usrs, nil
}

// UpdateUser updates a user in the store
func (s *Store) Update(ctx context.Context, traceID string, usr User) error {
	if usr.UID == "" {
		return fmt.Errorf("%s : missing UID", traceID)
	}

	if _, err := s.GetUserByUID(ctx, traceID, usr.UID); err != nil {
		return ErrNoExists
	}

	return s.update(ctx, traceID, usr)
}

// DeleteUser deletes a user from the store
func (s *Store) Delete(ctx context.Context, traceID string, usr User) error {
	if usr.UID == "" {
		return fmt.Errorf("%s : missing UID", traceID)
	}

	if _, err := s.GetUserByUID(ctx, traceID, usr.UID); err != nil {
		return ErrNoExists
	}

	return s.delete(ctx, traceID, usr.UID)
}

// ------ //

// add uses the client to add a user
func (s *Store) add(ctx context.Context, traceID string, usr User) (User, error) {
	jsonUser, err := json.Marshal(usr)
	if err != nil {
		return User{}, fmt.Errorf("unable to marshal user to json - %v", err)
	}

	mu := &api.Mutation{
		SetJson:   jsonUser,
		CommitNow: true,
	}

	s.log.Printf("%s : %s", traceID, "request to add user")

	resp, err := s.dgo.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return User{}, fmt.Errorf("unable to add user to db - %v", err)
	}

	if len(resp.Uids) == 0 {
		return User{}, fmt.Errorf("user id not returned - %v", resp.Json)
	}

	s.log.Printf("%s : %s : %s", traceID, "user added", usr.UID)

	usr.UID = resp.Uids["0"]
	return usr, nil
}

func (s *Store) query(ctx context.Context, traceID, q string, vars map[string]string) ([]User, error) {
	s.log.Printf("%s :request to query user", traceID)
	resp, err := s.dgo.NewTxn().QueryWithVars(ctx, q, vars)
	if err != nil {
		return []User{}, fmt.Errorf("dgo tx failed - QueryWithVars - %v", err)
	}

	type Result struct {
		Users []User `json:"query"`
	}

	var r Result
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return []User{}, fmt.Errorf("error while unmarshaling query result - %v", err)
	}

	if len(r.Users) < 1 {
		fmt.Println(len(r.Users))
		return []User{}, ErrNotFound
	}

	s.log.Printf("%s: returned %d users", traceID, len(r.Users))

	return r.Users, nil
}

func (s *Store) update(ctx context.Context, traceID string, usr User) error {
	mutation := &api.Mutation{
		CommitNow: true,
	}

	jsonUser, err := json.Marshal(usr)
	if err != nil {
		return fmt.Errorf("unable to marshal user to json - %v", err)
	}

	s.log.Printf("%s : %s", traceID, "request to update user")
	mutation.SetJson = jsonUser
	resp, err := s.dgo.NewTxn().Mutate(ctx, mutation)
	if err != nil {
		return fmt.Errorf("error updating user - %v", err)
	}

	if len(resp.Uids) != 1 {
		return fmt.Errorf("failed updating user\nReturned UIDs: %d", len(resp.Uids))
	}

	s.log.Printf("%s : %s : %s", traceID, "user updated successfully", usr.UID)

	return nil
}

func (s *Store) delete(ctx context.Context, traceID string, usrID string) error {
	mutation := &api.Mutation{
		CommitNow: true,
	}

	delUser := map[string]string{"uid": usrID}
	jsonUser, err := json.Marshal(delUser)
	if err != nil {
		return fmt.Errorf("%s : failed to marshal empty user - %v", traceID, err)
	}

	mutation.DeleteJson = jsonUser

	s.log.Printf("%s : %s", traceID, "request to delete user")
	_, err = s.dgo.NewTxn().Mutate(ctx, mutation)
	if err != nil {
		return fmt.Errorf("%s : unable to delete user - %v", traceID, err)
	}

	s.log.Printf("%s : %s : %s", traceID, "user deleted", usrID)

	return nil
}
