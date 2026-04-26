package infra

type InMemoryQueue struct {
	jobs chan string
}

func NewInMemoryQueue(buffer int) *InMemoryQueue {
	return &InMemoryQueue{jobs: make(chan string, buffer)}
}

func (q *InMemoryQueue) Enqueue(jobID string) {
	q.jobs <- jobID
}

func (q *InMemoryQueue) Dequeue() string {
	return <-q.jobs
}
