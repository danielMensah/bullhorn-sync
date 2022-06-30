package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
	"github.com/danielMensah/bullhorn-sync-poc/internal/repository"
)

func (c *Client) ConsumeCandidate(ctx context.Context, entity *bullhorn.Entity) error {
	switch entity.EventType {
	case bullhorn.EventTypeInserted:
		err := c.insert(ctx, entity)
		if err != nil {
			return err
		}
	case bullhorn.EventTypeUpdated:
		err := c.update(ctx, entity)
		if err != nil {
			return err
		}
	case bullhorn.EventTypeDeleted:
		err := c.delete(ctx, entity)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported event type: %s", entity.EventType)
	}

	return nil
}

func (c *Client) insert(ctx context.Context, entity *bullhorn.Entity) error {
	var candidate repository.Candidate

	err := json.Unmarshal(entity.Changes, &candidate)
	if err != nil {
		return err
	}

	// validations and transformations can happen here

	e := repository.DBEntity{
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

func (c *Client) update(ctx context.Context, entity *bullhorn.Entity) error {
	var candidate repository.Candidate

	err := json.Unmarshal(entity.Changes, &candidate)
	if err != nil {
		return err
	}

	// validations and transformations can happen here

	e := repository.DBEntity{
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

func (c *Client) delete(ctx context.Context, entity *bullhorn.Entity) error {
	// validations and transformations can happen here

	err := c.repo.Delete(ctx, entity.Id, entity.Name)
	if err != nil {
		return err
	}

	return nil
}
