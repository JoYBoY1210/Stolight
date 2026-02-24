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

func InitQueue(buffersize int) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	return &Queue{
		jobs:   make(chan string, buffersize),
		ctx:    ctx,
		cancel: cancel,
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
