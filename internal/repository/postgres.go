package repository

import (
	"context"
	"fmt"

	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresClient struct {
	db *pgxpool.Pool
}

func NewPostgresDB(ctx context.Context, cfg config.Config) (Repository, error) {
	conn, err := pgxpool.Connect(ctx, cfg.PostgresConnectionString)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to postgres: %w", err)
	}

	return &PostgresClient{
		db: conn,
	}, nil
}

func (p *PostgresClient) Save(ctx context.Context, dbEntity DBEntity) error {
	return nil
}

func (p *PostgresClient) Update(ctx context.Context, dbEntity DBEntity) error {
	return nil
}

func (p *PostgresClient) Delete(ctx context.Context, id int32, entityName string) error {
	return nil
}

func (p *PostgresClient) Close() {
	p.db.Close()
}
