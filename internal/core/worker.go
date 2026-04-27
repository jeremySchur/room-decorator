package core

import (
	"context"
	"log/slog"
	"room-decorator/internal/infra"
	"room-decorator/internal/models"
)

// RunWorker pulls job IDs off the queue and runs them. It does not yet
// support graceful shutdown: queue.Dequeue is uncancellable, so plumbing a
// ctx through here would be a lie. When Dequeue learns to honor cancellation,
// add a ctx parameter and pass it to the repo calls below.
func RunWorker(repo JobRepo, queue *infra.InMemoryQueue) {
	for {
		jobID := queue.Dequeue()

		job, err := repo.Get(context.Background(), jobID)
		if err != nil {
			slog.Warn("failed to load job, skipping", "job_id", jobID, "err", err)
			continue
		}

		if err := repo.UpdateStatus(context.Background(), jobID, models.Processing); err != nil {
			slog.Error("failed to mark job processing", "job_id", jobID, "err", err)
			continue
		}

		processErr := ProcessJob(job)

		nextStatus := models.Success
		if processErr != nil {
			nextStatus = models.Failed
		}
		if err := repo.UpdateStatus(context.Background(), jobID, nextStatus); err != nil {
			slog.Error("failed to update final job status", "job_id", jobID, "status", nextStatus, "err", err)
		}
	}
}
