package storage

import (
	"context"
	"fmt"

	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	pb "github.com/danielMensah/bullhorn-sync-poc/internal/proto"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
)

type Client struct {
	session gocqlx.Session
}

type Storage interface {
	Store(ctx context.Context, event *pb.Entity) error
	Update(ctx context.Context, event *pb.Entity) error
	Delete(ctx context.Context, entityName string, id int32) error
	Query(ctx context.Context, entityStruct interface{}, stmt string, names []string) error
	Close()
}

func New(config *config.Config) (Storage, error) {
	cluster := gocql.NewCluster(config.CassandraHosts...)
	cluster.Keyspace = config.CassandraKeyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.CassandraUsername,
		Password: config.CassandraPassword,
	}

	sess, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	return &Client{session: sess}, nil
}

func (c *Client) Query(ctx context.Context, entityStruct interface{}, stmt string, names []string) error {
	q := c.session.ContextQuery(ctx, stmt, names).BindStruct(entityStruct)
	if err := q.ExecRelease(); err != nil {
		return err
	}

	return nil
}

func (c *Client) Store(ctx context.Context, event *pb.Entity) error {
	return nil
}

func (c *Client) Update(ctx context.Context, event *pb.Entity) error {
	return nil
}

func (c *Client) Delete(ctx context.Context, entityName string, id int32) error {

}

func (c *Client) Close() {
	c.session.Close()
}
