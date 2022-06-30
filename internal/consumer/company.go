package consumer

import (
	"context"

	pb "github.com/danielMensah/bullhorn-sync-poc/internal/proto"
)

func (c *Consumer) ConsumeCompany(ctx context.Context, entity *pb.Entity) error {
	return nil
}
