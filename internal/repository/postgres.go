package repository

import (
	"context"
	"fmt"

	"github.com/danielMensah/bullhorn-sync-poc/internal/config"
	pb "github.com/danielMensah/bullhorn-sync-poc/internal/proto"
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

func (p *PostgresClient) Store(ctx context.Context, event *pb.Entity) error {
	return nil
}

func (p *PostgresClient) Update(ctx context.Context, event *pb.Entity) error {
	return nil
}

func (p *PostgresClient) Delete(ctx context.Context, entity *pb.Entity) error {
	return nil
}

func (p *PostgresClient) Close() {
	p.db.Close()
}
