package wpool

import (
	"context"

	"github.com/danielMensah/bullhorn-sync-poc/internal/bullhorn"
)

type JobID string
type jobType string
type jobMetadata map[string]interface{}

type ExecutionFn func(context.Context, *bullhorn.Entity) error

type JobDescriptor struct {
	ID       JobID
	JType    jobType
	Metadata jobMetadata
}

type Result struct {
	Err        error
	Descriptor JobDescriptor
}

type Job struct {
	Descriptor JobDescriptor
	ExecFn     ExecutionFn
	Entity     *bullhorn.Entity
}

func (j Job) execute(ctx context.Context) Result {
	err := j.ExecFn(ctx, j.Entity)
	if err != nil {
		return Result{
			Err:        err,
			Descriptor: j.Descriptor,
		}
	}

	return Result{
		Descriptor: j.Descriptor,
	}
}
