package models

import "time"

type JobStatus string

const (
	Pending    JobStatus = "PENDING"
	Processing JobStatus = "PROCESSING"
	Success    JobStatus = "SUCCESS"
	Failed     JobStatus = "FAILED"
)

type Job struct {
	ID        string    `json:"id"`
	Status    JobStatus `json:"status"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
