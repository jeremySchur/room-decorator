package core

import (
	"log"
	"room-decorator/internal/infra"
	"room-decorator/internal/models"
)

func RunWorker(repo *infra.InMemoryJobRepo, queue *infra.InMemoryQueue) {
	for {
		jobID := queue.Dequeue()

		job, ok := repo.Get(jobID)
		if !ok {
			log.Printf("job %s not found, skipping", jobID)
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
