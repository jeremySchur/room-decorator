package core

import (
	"context"
	"log/slog"
	"room-decorator/internal/infra"
	"room-decorator/internal/models"
	"time"

	"github.com/google/uuid"
)

func CreateJob(ctx context.Context, repo JobRepo, queue *infra.InMemoryQueue, payload string) (*models.Job, error) {
	now := time.Now().UTC()
	job := &models.Job{
		ID:        uuid.NewString(),
		Status:    models.Pending,
		Payload:   payload,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := repo.Insert(ctx, job); err != nil {
		return nil, err
	}
	queue.Enqueue(job.ID)
	return job, nil
}

func ProcessJob(job *models.Job) error {
	slog.Info("processing job", "job_id", job.ID)
	time.Sleep(500 * time.Millisecond)
	slog.Info("finished processing job", "job_id", job.ID)

	return nil
}
