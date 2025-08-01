package schema

import (
	"bytes"
	"context"
	"dgraph-client/data/role"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"time"

	"github.com/charmbracelet/log"
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
)

//go:embed schema.dgraph
var schemaDoc string

// errors
var (
	ErrNoSchemaFound = errors.New("no shchema found")
	ErrInvalidSchema = errors.New("invalid schema")
)

// schema consists of a schema doc and dgraph client
type Schema struct {
	dgo    *dgo.Dgraph
	schema string
}

// NewSchema initializes our Schema type
func NewSchema(dgo *dgo.Dgraph) (*Schema, error) {
	tmpl := template.New("schema")
	if _, err := tmpl.Parse(schemaDoc); err != nil {
		return nil, fmt.Errorf("schema - template parse error - %v", err)
	}

	var schemaBuff bytes.Buffer
	if err := tmpl.Execute(&schemaBuff, nil); err != nil {
		return nil, fmt.Errorf("schema - unable to execute template - %v", err)
	}
	schema := &Schema{
		dgo:    dgo,
		schema: schemaBuff.String(),
	}

	return schema, nil
}

// InitSchema creates a schema in our database
func (s *Schema) InitSchema(ctx context.Context) error {
	op := &api.Operation{}
	op.Schema = s.schema

	if err := s.dgo.Alter(ctx, op); err != nil {
		return fmt.Errorf("schema - InitDB error - %v", err)
	}
	return nil
}

// InitRoles creates the default roles in our database
func (s *Schema) InitRoles(ctx context.Context, log *log.Logger, traceID string) error {
	rs := role.NewStore(log, s.dgo)
	roles := []string{"admin", "user"}

	txn := s.dgo.NewTxn()
	defer txn.Discard(ctx)

	for _, r := range roles {
		fmt.Println("Role: ", r)
		role, err := rs.Add(ctx, traceID, r, time.Now())
		if err != nil {
			return fmt.Errorf("unable to add new role - %s", err)
		}

		jsonRole, err := json.Marshal(role)
		if err != nil {
			return fmt.Errorf("unable to marshal admin role to json - %v", err)
		}

		mu := &api.Mutation{
			SetJson: jsonRole,
		}

		fmt.Println("adding roles")

		resp, err := txn.Mutate(ctx, mu)
		if err != nil {
			return fmt.Errorf("unable to add role to db - %v", err)
		}

		if len(resp.Uids) == 0 {
			return fmt.Errorf("role id not returned - %v", resp.Json)
		}
	}

	err := txn.Commit(ctx)
	if err != nil {
		return fmt.Errorf("unable to commit transaction - %v", err)
	}

	return nil
}

// DropData drops all data from the database but leaves the schema
func (s *Schema) DropData(ctx context.Context) error {
	if err := s.dgo.Alter(ctx, &api.Operation{DropOp: api.Operation_DATA}); err != nil {
		return fmt.Errorf("schema - dropData error - %s", err)
	}

	return nil
}

// DropAll drops all data and the schema - clean slate
func (s *Schema) DropAll(ctx context.Context) error {
	// clear the db for the example
	if err := s.dgo.Alter(ctx, &api.Operation{DropOp: api.Operation_ALL}); err != nil {
		return fmt.Errorf("schema - dropAll error - %s", err)
	}

	return nil
}
