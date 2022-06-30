package repository

import (
	"context"
)

type DBEntity struct {
	ID   int32
	Name string
	Data interface{}
}

type Repository interface {
	Save(ctx context.Context, dbEntity DBEntity) error
	Update(ctx context.Context, dbEntity DBEntity) error
	Delete(ctx context.Context, id int32, entityName string) error
	Close()
}
