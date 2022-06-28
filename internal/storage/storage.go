package storage

import (
	"context"
	"fmt"

	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	pb "github.com/danielMensah/bullhorn-sync-poc/internal/proto"
	"github.com/gocql/gocql"
)

type Client struct {
	svc *gocql.Session
}

type Storage interface {
	Store(ctx context.Context, event *pb.Event) error
	Update(ctx context.Context, event *pb.Event) error
	Delete(ctx context.Context, id int32) error
	Close()
}

func New(config *config.Config) (Storage, error) {
	cluster := gocql.NewCluster(config.CassandraHosts...)
	cluster.Keyspace = config.CassandraKeyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.CassandraUsername,
		Password: config.CassandraPassword,
	}

	sess, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	return &Client{svc: sess}, nil
}

func (c *Client) Store(ctx context.Context, event *pb.Event) error {
	return nil
}

func (c *Client) Update(ctx context.Context, event *pb.Event) error {
	return nil
}

func (c *Client) Delete(ctx context.Context, id int32) error {
	return nil
}

func (c *Client) Close() {
	c.svc.Close()
}
