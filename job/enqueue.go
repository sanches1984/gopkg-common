package job

import (
	"context"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type Enqueue struct {
	enqueuer *work.Enqueuer
}

func NewEnqueue(ns string, pool *redis.Pool) *Enqueue {
	return &Enqueue{
		enqueuer: work.NewEnqueuer(ns, pool),
	}
}

func (e *Enqueue) AddJob(ctx context.Context, job interface{}) error {
	jobName, err := getJobName(ctx, job)
	if err != nil {
		return err
	}
	args, err := packArguments(ctx, job)
	if err != nil {
		return err
	}
	_, err = e.enqueuer.EnqueueUnique(jobName, args)
	return err
}
