package queue

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type Queue struct {
	jobs   chan string
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

var queue *Queue

func InitQueue(buffersize int) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	q := &Queue{
		jobs:   make(chan string, buffersize),
		ctx:    ctx,
		cancel: cancel,
	}
	queue = q
	return q
}

func GetQueue() *Queue {
	return queue
}	

func (q *Queue) StartWorkerPool(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		q.wg.Add(1)
		go func(workerID int) {
			defer q.wg.Done()
			for {
				select {
				case <-q.ctx.Done():

					log.Printf("Worker %d stopping gracefully..", workerID)
					return

				case id := <-q.jobs:

					err := Worker(id)
					if err != nil {
						log.Printf("Worker %d error processing %s: %v", workerID, id, err)
					}
				}
			}
		}(i)
	}
}

func (q *Queue) AddJob(jobID string) error {
	select {
	case q.jobs <- jobID:
		return nil
	case <-q.ctx.Done():
		return fmt.Errorf("Queue is closed")

	}
}

func (q *Queue) Close() {
	q.cancel()
	q.wg.Wait()
	log.Println("Queue closed and all workers stopped")
}
