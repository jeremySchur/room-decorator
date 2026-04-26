package infra

import "testing"

func TestQueueEnqueueDequeue(t *testing.T) {
	q := NewInMemoryQueue(2)
	q.Enqueue("a")
	q.Enqueue("b")
	if got := q.Dequeue(); got != "a" {
		t.Errorf("want a, got %q", got)
	}
	if got := q.Dequeue(); got != "b" {
		t.Errorf("want b, got %q", got)
	}
}
