package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/mikeheft/go-backend/db/sqlc"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
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
	logger := NewLogger()
	redis.SetLogger(logger)

	server := asynq.NewServer(
		redisClientOpt,
		asynq.Config{
			Queues: map[string]int{
				CriticalQueue: 10,
				DefaultQueue:  2,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				// This can also handle sending email, jiraTicket, or slack notificiations
				log.Error().
					Err(err).
					Str("type", task.Type()).
					Bytes("payload", task.Payload()).
					Msg("process task failed")
			}),
			Logger: logger,
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
