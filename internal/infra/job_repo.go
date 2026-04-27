package infra

import (
	"context"
	"errors"
	"room-decorator/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresJobRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresJobRepo(pool *pgxpool.Pool) *PostgresJobRepo {
	return &PostgresJobRepo{pool: pool}
}

func (r *PostgresJobRepo) Get(ctx context.Context, id string) (*models.Job, error) {
	const q = `SELECT id, status, payload, created_at, updated_at FROM jobs WHERE id = $1`

	var job models.Job
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&job.ID, &job.Status, &job.Payload, &job.CreatedAt, &job.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, models.ErrJobNotFound
	}
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *PostgresJobRepo) Insert(ctx context.Context, job *models.Job) error {
	const q = `INSERT INTO jobs (id, status, payload, created_at, updated_at)
	           VALUES ($1, $2, $3, $4, $5)`

	_, err := r.pool.Exec(ctx, q, job.ID, job.Status, job.Payload, job.CreatedAt, job.UpdatedAt)
	return err
}

func (r *PostgresJobRepo) UpdateStatus(ctx context.Context, id string, status models.JobStatus) error {
	const q = `UPDATE jobs SET status = $1, updated_at = $2 WHERE id = $3`

	_, err := r.pool.Exec(ctx, q, status, time.Now().UTC(), id)
	return err
}
