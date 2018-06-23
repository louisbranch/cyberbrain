package test

import (
	"errors"

	"gitlab.com/luizbranco/cyberbrain/primitives"
)

type WorkerPool struct {
	RegisterFunc func(string, primitives.Worker) error
	EnqueueFunc  func(string, interface{}) error
}

func (p *WorkerPool) Register(name string, worker primitives.Worker) error {
	if p.RegisterFunc == nil {
		return errors.New("RegisterFunc not implemented")
	}

	return p.RegisterFunc(name, worker)
}

func (p *WorkerPool) Enqueue(name string, v interface{}) error {
	if p.EnqueueFunc == nil {
		return errors.New("EnqueueFunc not implemented")
	}

	return p.EnqueueFunc(name, v)
}
