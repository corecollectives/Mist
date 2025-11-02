package queue

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Queue struct {
	jobs chan int64

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewQueue(buffer int) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	q := &Queue{
		jobs: make(chan int64, buffer),

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
		for id := range q.jobs {

			fmt.Printf("job with job id %d has started\n", id)
			time.Sleep(2 * time.Second)
			fmt.Printf("job with job id %d has finished\n", id)

		}

	}()
}

func (q *Queue) AddJob(Id int64) error {
	select {
	case q.jobs <- Id:
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
