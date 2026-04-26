package api

import (
	"encoding/json"
	"log"
	"net/http"
	"room-decorator/internal/core"
	"room-decorator/internal/infra"
)

type Server struct {
	repo  *infra.InMemoryJobRepo
	queue *infra.InMemoryQueue
}

func NewServer(repo *infra.InMemoryJobRepo, queue *infra.InMemoryQueue) *Server {
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

	job := core.CreateJob(s.repo, s.queue, req.Payload)

	writeJSON(w, http.StatusCreated, job)
}

func (s *Server) getJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	job, ok := s.repo.Get(id)
	if !ok {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, job)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("failed to write JSON response: %v", err)
	}
}
