package queue

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID         int
	AppID      int
	CommitHash string
	Logs       string
	Status     string
	CreatedAt  time.Time
	FinishedAt time.Time
}

type Queue struct {
	jobs chan Job

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewQueue(buffer int) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	q := &Queue{
		jobs: make(chan Job, buffer),

		ctx:    ctx,
		cancel: cancel,
	}
	q.StartWorker()
	return q

}

func (q *Queue) StartWorker() {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		for job := range q.jobs {

			fmt.Printf("job with job id %d has started\n", job.ID)
			time.Sleep(2 * time.Second)
			fmt.Printf("job with job id %d has finished\n", job.ID)

		}

	}()
}

func (q *Queue) AddJob(job Job) error {
	select {
	case q.jobs <- job:
		return nil
	case <-q.ctx.Done():
		return fmt.Errorf("queue is closed")
	default:
		return fmt.Errorf("queue is full")
	}
}

func (q *Queue) Close() {
	q.cancel()
	close(q.jobs)
	q.wg.Wait()
	fmt.Println("queue is closed")
}
