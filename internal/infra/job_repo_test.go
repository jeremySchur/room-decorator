package infra

import (
	"room-decorator/internal/models"
	"testing"
	"time"
)

func newTestJob(id string) *models.Job {
	now := time.Now().UTC()
	return &models.Job{
		ID:        id,
		Status:    models.Pending,
		Payload:   "test-payload",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestInsertThenGet_ReturnsSameJob(t *testing.T) {
	repo := NewInMemoryJobRepo()
	job := newTestJob("1")

	repo.Insert(job)

	got, ok := repo.Get(job.ID)
	if !ok {
		t.Fatalf("Get(%q) returned ok=false, want true", job.ID)
	}
	if got != job {
		t.Errorf("Get(%q) = %+v, want %+v", job.ID, got, job)
	}
}

func TestGet_UnknownID_ReturnsNotFound(t *testing.T) {
	repo := NewInMemoryJobRepo()

	got, ok := repo.Get("does-not-exist")
	if ok {
		t.Errorf("Get returned ok=true for unknown ID, want false")
	}
	if got != nil {
		t.Errorf("Get returned %+v for unknown ID, want nil", got)
	}
}

func TestUpdateStatus(t *testing.T) {
	repo := NewInMemoryJobRepo()
	job := newTestJob("1")
	repo.Insert(job)

	repo.UpdateStatus(job.ID, models.Processing)

	got, _ := repo.Get(job.ID)
	if got.Status != models.Processing {
		t.Errorf("Status = %v, want %v", got.Status, models.Processing)
	}
}

func TestUpdateStatus_BumpsUpdatedAt(t *testing.T) {
	repo := NewInMemoryJobRepo()
	job := newTestJob("1")
	originalUpdatedAt := job.UpdatedAt
	repo.Insert(job)

	time.Sleep(time.Millisecond) // ensure clock moved forward
	repo.UpdateStatus(job.ID, models.Processing)

	got, _ := repo.Get(job.ID)
	if !got.UpdatedAt.After(originalUpdatedAt) {
		t.Errorf("UpdatedAt = %v, want it to be after %v", got.UpdatedAt, originalUpdatedAt)
	}
}

func TestUpdatesStatus_UnknownId_DoesNotPanic(t *testing.T) {
	repo := NewInMemoryJobRepo()

	repo.UpdateStatus("does-not-exist", models.Failed) // should be a no-op
}
