package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/danielMensah/bullhorn-sync-poc/internal/proto"
	"github.com/danielMensah/bullhorn-sync-poc/internal/repository"
)

func (c *Consumer) ConsumeCandidate(ctx context.Context, entity *pb.Entity) error {
	switch entity.EventType {
	case pb.EventType_INSERTED:
		err := c.insertCandidate(ctx, entity)
		if err != nil {
			return err
		}
	case pb.EventType_UPDATED:
		err := c.updateCandidate(ctx, entity)
		if err != nil {
			return err
		}
	case pb.EventType_DELETED:
		err := c.deleteCandidate(ctx, entity)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported event type: %s", entity.EventType)
	}

	return nil
}

func (c *Consumer) insertCandidate(ctx context.Context, entity *pb.Entity) error {
	var candidate repository.Candidate

	err := json.Unmarshal(entity.Changes, &candidate)
	if err != nil {
		return err
	}

	// validations and transformations can happen here

	e := repository.Entity{
		ID:   entity.Id,
		Name: entity.Name,
		Data: candidate,
	}

	err = c.repo.Save(ctx, e)
	if err != nil {
		return err
	}

	return nil
}

func (c *Consumer) updateCandidate(ctx context.Context, entity *pb.Entity) error {
	var candidate repository.Candidate

	err := json.Unmarshal(entity.Changes, &candidate)
	if err != nil {
		return err
	}

	// validations and transformations can happen here

	e := repository.Entity{
		ID:   entity.Id,
		Name: entity.Name,
		Data: candidate,
	}

	err = c.repo.Update(ctx, e)
	if err != nil {
		return err
	}

	return nil
}

func (c *Consumer) deleteCandidate(ctx context.Context, entity *pb.Entity) error {
	// validations and transformations can happen here

	err := c.repo.Delete(ctx, entity.Id, entity.Name)
	if err != nil {
		return err
	}

	return nil
}
