package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	db "github.com/leedrum/simplebank/db/sqlc"
	"github.com/leedrum/simplebank/util"
	"github.com/rs/zerolog/log"
)

const TaskTypeSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TaskTypeSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return err
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msgf("task enqueued: id=%s type=%s", info.ID, TaskTypeSendVerifyEmail)
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(
	ctx context.Context,
	task *asynq.Task,
) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return fmt.Errorf("user not found: %w", asynq.SkipRetry)
		// }

		return fmt.Errorf("failed to get user: %w", err)
	}

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})

	subject := "Verify your email"
	body := fmt.Sprintf(
		"Hello, %s! Click this link to verify your email: http://localhost:8080/v1/verify_email?email_id=%v&code=%s",
		verifyEmail.Username,
		verifyEmail.ID,
		verifyEmail.SecretCode,
	)
	to := []string{user.Email}
	processor.mailer.SendEmail(
		subject,
		body,
		to,
		nil, nil, nil,
	)

	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Msgf("processing task: id=%s type=%s", task.ResultWriter().TaskID(), TaskTypeSendVerifyEmail)

	return nil
}
