package worker

import (
	"fmt"

	"github.com/sofiukl/oms/oms-checkout/api"
)

// Worker - This is Worker struct
type Worker struct {
	ID          int
	WorkChannel chan Work
	WorkerQueue chan chan Work
}

// NewWorker - This creates the instance of new worker
func NewWorker(id int, workerQueue chan chan Work) Worker {
	worker := Worker{
		ID:          id,
		WorkChannel: make(chan Work),
		WorkerQueue: workerQueue,
	}

	return worker
}

// Start - This is the runnable method of the worker
func (w *Worker) Start() {
	go func() {
		for {
			// assigning available channel to WorkerQueue
			w.WorkerQueue <- w.WorkChannel
			select {
			case job := <-w.WorkChannel:

				// Receive a work request.
				fmt.Printf("worker %d: working on %s!\n", w.ID, job.Work.CartID)
				//time.Sleep(40000 * time.Millisecond)
				api.CheckoutProduct(job.Conn, job.Config, job.Work, job.Lock)
				fmt.Printf("worker%d: work %s completed\n", w.ID, job.Work.CartID)
			}
		}
	}()
}
