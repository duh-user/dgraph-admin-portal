// Package data contains the schema, models, and client
// to access our dgraph instance
package data

import (
	"context"
	"dgraph-client/config"
	"fmt"
	"log"
	"time"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
)

type CancelFunc func()

type DGClient struct {
	Client *dgo.Dgraph
}

func NewDGClient(cfg *config.Config) (DGClient, CancelFunc) {
	// TODO - remove withinsecure for production
	// TODO - add TLS to encrypt grpc connection
	conn, err := grpc.Dial(cfg.DGAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("gRPC dial error - %w", cfg.DGAddr, "-", err)
	}

	fmt.Println(cfg.DGAddr)
	client := dgo.NewDgraphClient(api.NewDgraphClient(conn))
	dgclient := DGClient{
		Client: client,
	}

	return dgclient, func() {
		if err := conn.Close(); err != nil {
			log.Println("dgraph server conn error - %w", err)
		}
	}
}

func (dgc *DGClient) HealthCheck(ctx context.Context, retryInterval time.Duration) error {
	var t *time.Timer

	for {
		if err := healthCheck(ctx, dgc.Client); err == nil {
			return nil
		}

		if ctx.Err() != nil {
			return fmt.Errorf("health check timed out - %w", ctx.Err())
		}

		if t == nil {
			t = time.NewTimer(retryInterval)
		}

		select {
		case <-ctx.Done():
			t.Stop()
			return fmt.Errorf("health check timed out - %w", ctx.Err())
		case <-t.C:
			t.Reset(retryInterval)
		}
	}
}

func healthCheck(ctx context.Context, client *dgo.Dgraph) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	_, err := client.NewTxn().Query(ctx, `{
			result(func: has(uid), first: 1) {
				uid
			}
		}
	}`)
	if err != nil {
		return err
	}
	return nil
}
