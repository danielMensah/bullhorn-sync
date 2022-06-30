package repository

import (
	"context"
	"fmt"

	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/table"
)

var (
	mappingTable = map[string]*table.Table{
		"candidate": CandidateTable,
	}
)

type CassandraClient struct {
	db gocqlx.Session
}

func NewCassandraDB(config *config.Config) (Repository, error) {
	cluster := gocql.NewCluster(config.CassandraHosts...)
	cluster.Keyspace = config.CassandraKeyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: config.CassandraUsername,
		Password: config.CassandraPassword,
	}

	sess, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %v", err)
	}

	return &CassandraClient{db: sess}, nil
}

func (c *CassandraClient) Save(ctx context.Context, entityStruct Entity) error {
	entityTable, ok := mappingTable[entityStruct.Name]
	if !ok {
		return fmt.Errorf("entity '%s' not supported", entityStruct.Name)
	}

	q := entityTable.InsertQueryContext(ctx, c.db).BindStruct(entityStruct.Data)
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("cannot exec save query: %w", err)
	}

	return nil
}

func (c *CassandraClient) Update(ctx context.Context, entityStruct Entity) error {
	entityTable, ok := mappingTable[entityStruct.Name]
	if !ok {
		return fmt.Errorf("entity '%s' not supported", entityStruct.Name)
	}

	q := entityTable.UpdateQueryContext(ctx, c.db).BindStruct(entityStruct)
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("cannot exec update query: %w", err)
	}

	return nil
}

func (c *CassandraClient) Delete(ctx context.Context, id int32, entityName string) error {
	entityTable, ok := mappingTable[entityName]
	if !ok {
		return fmt.Errorf("entity '%s' not supported", entityName)
	}

	q := entityTable.DeleteQueryContext(ctx, c.db, "id").Bind(id)
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("cannot exec update query: %w", err)
	}

	return nil
}

func (c *CassandraClient) Close() {
	c.db.Close()
}
