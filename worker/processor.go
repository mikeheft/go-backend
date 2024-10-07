package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/mikeheft/go-backend/db/sqlc"
)

const (
	CriticalQueue = "critical"
	DefaultQueue  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessSendVerifyEmailTask(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcesser struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisClientOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(
		redisClientOpt,
		asynq.Config{
			Queues: map[string]int{
				CriticalQueue: 10,
				DefaultQueue:  2,
			},
		},
	)

	return &RedisTaskProcesser{
		server: server,
		store:  store,
	}
}

func (processor *RedisTaskProcesser) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessSendVerifyEmailTask)

	return processor.server.Start(mux)
}
