package worker

import (
	"errors"
	"fmt"
	"time"

	"gitlab.com/luizbranco/srs/primitives"
)

type Worker struct {
	Database primitives.Database
}

func (w *Worker) Start() error {
	tick := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-tick.C:
			fmt.Printf(".")
		}
	}

	return nil
}

func (w *Worker) Register(name string, worker primitives.Worker) error {
	return errors.New("not implemented")
}

func (w *Worker) Enqueue() error {
	return errors.New("not implemented")
}
