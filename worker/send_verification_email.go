package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type VerificationEmailPayload struct {
	Username string `json:"username"`
}

const (
	TaskSendVerifyEmail = "task:send_verify_email"
)

func (distributor *RedisTaskDistributor) SendVerificationEmailTask(
	ctx context.Context,
	payload VerificationEmailPayload,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", jsonPayload).
		Str("queue", info.Queue).
		Int("maxy_retry", info.MaxRetry).
		Msg("enqued task")

	return nil
}
