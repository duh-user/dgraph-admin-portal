package role

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
)

// Errors
var (
	ErrNoExists     = errors.New("role does not exit")
	ErrExists       = errors.New("role exists")
	ErrNotFound     = errors.New("role not found")
	ErrPassNotMatch = errors.New("passwords do not match")
)

// Store will manage the role store API's
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

// Add will add a new role to the db if the user doesn't already exist
// if the role existss the found user is returned
// if added the role with uid is returned
func (s *Store) Add(ctx context.Context, traceID string, role string, now time.Time) (Role, error) {
	if r, err := s.GetRoleByName(ctx, traceID, role); err == nil {
		return r, ErrExists
	}

	r := Role{
		Name:         role,
		DateCreated:  now,
		LastSeen:     now,
		LastModified: now,
	}

	return s.add(ctx, traceID, r)
}

func (s *Store) GetRoleByName(ctx context.Context, traceID string, name string) (Role, error) {
	vars := make(map[string]string)
	vars["$role_name"] = name
	// all queries need to start with the name "query" to work with our query handler
	q := `
			query query($role_name: string){
				query(func: eq(role_name, $role_name))	{
					uid
					role_name
					date_created
					last_modified
					last_seen
				}	
			}	
		`

	role, err := s.query(ctx, traceID, q, vars)
	if err != nil {
		return Role{}, err
	}

	return role[0], nil
}

// --- Internal Functions

func (s *Store) add(ctx context.Context, traceID string, role Role) (Role, error) {
	jsonRole, err := json.Marshal(role)
	if err != nil {
		return Role{}, fmt.Errorf("unable to marshal role to json - %v", err)
	}

	mu := &api.Mutation{
		SetJson:   jsonRole,
		CommitNow: true,
	}

	s.log.Printf("%s : %s", traceID, "request to add role")

	resp, err := s.dgo.NewTxn().Mutate(ctx, mu)
	if err != nil {
		return Role{}, fmt.Errorf("unable to add role to db - %v", err)
	}

	if len(resp.Uids) == 0 {
		return Role{}, fmt.Errorf("role uid not returned - %v", resp.Json)
	}

	s.log.Printf("%s : %s : %s", traceID, "role added", role.UID)

	role.UID = resp.Uids["0"]
	return role, nil
}

func (s *Store) query(ctx context.Context, traceID, q string, vars map[string]string) ([]Role, error) {
	s.log.Printf("%s : %s", traceID, "request to query role")
	resp, err := s.dgo.NewTxn().QueryWithVars(ctx, q, vars)
	if err != nil {
		return []Role{}, fmt.Errorf("dgo tx failed - QueryWithVars - %v", err)
	}

	type Result struct {
		Roles []Role `json:"query"`
	}
	var r Result
	//fmt.Println(string(resp.Json))
	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		return []Role{}, fmt.Errorf("error while unmarshaling query result - %v", err)
	}

	if len(r.Roles) < 1 {
		fmt.Println(len(r.Roles))
		return []Role{}, ErrNotFound
	}

	s.log.Printf("%s: %d %s", traceID, len(r.Roles), "roles returned")

	return r.Roles, nil
}

func (s *Store) update(ctx context.Context, traceID string, usr Role) error {
	mutation := &api.Mutation{
		CommitNow: true,
	}

	jsonRole, err := json.Marshal(usr)
	if err != nil {
		return fmt.Errorf("unable to marshal role to json - %v", err)
	}

	s.log.Printf("%s : %s", traceID, "request to update role")
	mutation.SetJson = jsonRole
	resp, err := s.dgo.NewTxn().Mutate(ctx, mutation)
	if err != nil {
		return fmt.Errorf("error updating role - %v", err)
	}

	if len(resp.Uids) != 1 {
		return fmt.Errorf("failed updating role\nReturned UIDs: %d", len(resp.Uids))
	}

	s.log.Printf("%s : %s : %s", traceID, "role updated successfully", usr.UID)

	return nil
}

func (s *Store) delete(ctx context.Context, traceID string, usrID string) error {
	mutation := &api.Mutation{
		CommitNow: true,
	}

	delRole := map[string]string{"uid": usrID}
	jsonRole, err := json.Marshal(delRole)
	if err != nil {
		return fmt.Errorf("%s : failed to marshal empty role - %v", traceID, err)
	}

	mutation.DeleteJson = jsonRole

	s.log.Printf("%s : %s", traceID, "request to delete role")
	_, err = s.dgo.NewTxn().Mutate(ctx, mutation)
	if err != nil {
		return fmt.Errorf("%s : unable to delete role - %v", traceID, err)
	}

	s.log.Printf("%s : %s : %s", traceID, "role deleted", usrID)

	return nil
}
