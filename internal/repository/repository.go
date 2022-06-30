package repository

import (
	"context"
)

type Entity struct {
	ID   int32
	Name string
	Data interface{}
}

type Repository interface {
	Save(ctx context.Context, entityStruct Entity) error
	Update(ctx context.Context, entityStruct Entity) error
	Delete(ctx context.Context, id int32, entityName string) error
	Close()
}
