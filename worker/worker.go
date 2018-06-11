package worker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/luizbranco/srs/primitives"
)

type WorkerPool struct {
	workers map[string]primitives.Worker

	Database primitives.Database
}

func (w *WorkerPool) Start() error {
	tick := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-tick.C:
			fmt.Printf(".")
		}
	}

	return nil
}

func (p *WorkerPool) Register(name string, worker primitives.Worker) error {
	if p.workers == nil {
		p.workers = make(map[string]primitives.Worker)
	}

	_, ok := p.workers[name]
	if ok {
		return errors.Errorf("%s already registered", name)
	}

	p.workers[name] = worker

	return nil
}

func (w *WorkerPool) Enqueue(name string, args map[string]string) error {
	b, err := json.Marshal(args)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal job args %q", name)
	}

	job := &Job{
		Name: name,
		Args: b,
	}

	err = w.Database.Create(job)
	if err != nil {
		return errors.Wrapf(err, "failed to save job to database %q", name)
	}

	return nil
}
