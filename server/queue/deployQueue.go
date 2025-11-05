package queue

import (
	"context"
	"fmt"
	"sync"

	"github.com/corecollectives/mist/api/handlers/dockerdeploy"
)

type Queue struct {
	jobs chan int64

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewQueue(buffer int, d *dockerdeploy.Deployer) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	q := &Queue{
		jobs: make(chan int64, buffer),

		ctx:    ctx,
		cancel: cancel,
	}
	q.StartWorker(d)
	return q

}

func (q *Queue) StartWorker(d *dockerdeploy.Deployer) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		for id := range q.jobs {

			d.DeployerMain(id)

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
