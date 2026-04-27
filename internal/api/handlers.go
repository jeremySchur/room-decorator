package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"room-decorator/internal/core"
	"room-decorator/internal/infra"
	"room-decorator/internal/models"
)

type Server struct {
	repo  core.JobRepo
	queue *infra.InMemoryQueue
}

func NewServer(repo core.JobRepo, queue *infra.InMemoryQueue) *Server {
	return &Server{repo: repo, queue: queue}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /jobs", s.createJob)
	mux.HandleFunc("GET /jobs/{id}", s.getJob)
	return mux
}

type createJobRequest struct {
	Payload string `json:"payload"`
}

func (s *Server) createJob(w http.ResponseWriter, r *http.Request) {
	var req createJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if req.Payload == "" {
		http.Error(w, "payload is required", http.StatusBadRequest)
		return
	}

	job, err := core.CreateJob(r.Context(), s.repo, s.queue, req.Payload)
	if err != nil {
		slog.Error("create job failed", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, job)
}

func (s *Server) getJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	job, err := s.repo.Get(r.Context(), id)
	if errors.Is(err, models.ErrJobNotFound) {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("get job failed", "err", err, "job_id", id)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, job)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("failed to write JSON response", "err", err)
	}
}
