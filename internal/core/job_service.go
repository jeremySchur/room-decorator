package core

import (
	"log"
	"room-decorator/internal/infra"
	"room-decorator/internal/models"
	"time"

	"github.com/google/uuid"
)

func CreateJob(repo *infra.InMemoryJobRepo, queue *infra.InMemoryQueue, payload string) *models.Job {
	job := &models.Job{
		ID:        uuid.NewString(),
		Status:    models.Pending,
		Payload:   payload,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	repo.Insert(job)
	queue.Enqueue(job.ID)
	return job
}

func ProcessJob(job *models.Job) error {
	log.Printf("processing job %s", job.ID)
	time.Sleep(500 * time.Millisecond)
	log.Printf("finished processing job %s", job.ID)

	return nil
}
