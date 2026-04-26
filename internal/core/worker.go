package core

import (
	"log/slog"
	"room-decorator/internal/infra"
	"room-decorator/internal/models"
)

func RunWorker(repo *infra.InMemoryJobRepo, queue *infra.InMemoryQueue) {
	for {
		jobID := queue.Dequeue()

		job, ok := repo.Get(jobID)
		if !ok {
			slog.Warn("job not found, skipping", "job_id", jobID)
			continue
		}

		repo.UpdateStatus(jobID, models.Processing)

		err := ProcessJob(job)

		if err != nil {
			repo.UpdateStatus(jobID, models.Failed)
		} else {
			repo.UpdateStatus(jobID, models.Success)
		}
	}
}
