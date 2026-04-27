package core

import (
	"context"
	"room-decorator/internal/models"
)

type JobRepo interface {
	Get(ctx context.Context, id string) (*models.Job, error)
	Insert(ctx context.Context, job *models.Job) error
	UpdateStatus(ctx context.Context, id string, status models.JobStatus) error
}
