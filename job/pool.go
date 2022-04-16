package job

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

func NewPool(ctx context.Context, concurrency uint, redisNS string, redisPool *redis.Pool, jobList []WorkRecord, jobCtx interface{}, middlewareList []interface{}, logger zerolog.Logger) *work.WorkerPool {
	pool := work.NewWorkerPool(jobCtx, uint(concurrency), redisNS, redisPool)
	for _, mw := range middlewareList {
		pool.Middleware(mw)
	}
	for _, item := range jobList {
		jobName, err := getJobName(ctx, item.Job)
		if err != nil {
			panic(fmt.Sprintf("can't get job name: %v", err))
		}

		if item.Schedule != "" {
			pool.Job(jobName, item.Fn)
			pool.PeriodicallyEnqueue(item.Schedule, jobName)

			logger.Info().Str("job", jobName).Str("schedule", item.Schedule).Msg("add schedule for job")
		} else {
			pool.Job(jobName, item.Fn)

			logger.Info().Str("job", jobName).Msg("add job")
		}
	}
	return pool
}
