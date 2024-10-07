package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	SendVerificationEmailTask(
		ctx context.Context,
		payload *VerificationEmailPayload,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisClientOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisClientOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
