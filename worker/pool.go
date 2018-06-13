package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/luizbranco/cyberbrain/primitives"
)

type WorkerPool struct {
	workers map[string]primitives.Worker

	Database primitives.Database
}

func (w *WorkerPool) Start() {
	tick := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-tick.C:
			w.run()
		}
	}
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
		Name:  name,
		Args:  b,
		State: scheduled,
		RunAt: time.Now(),
	}

	err = w.Database.Create(job)
	if err != nil {
		return errors.Wrapf(err, "failed to save job to database %q", name)
	}

	return nil
}

func (wp *WorkerPool) run() {
	q := &query{}

	rs, err := wp.Database.Query(q)
	if err != nil {
		err = errors.Wrap(err, "failed to query scheduled jobs")
		log.Println(err)
		return
	}

	var jobs []Job

	for _, r := range rs {
		job, ok := r.(*Job)
		if !ok {
			err := errors.Errorf("invalid record type %T", r)
			log.Println(err)
			return
		}

		jobs = append(jobs, *job)
	}

	for _, j := range jobs {
		wp.runJob(j)
	}
}

func (wp *WorkerPool) runJob(j Job) {
	worker, ok := wp.workers[j.Name]

	if !ok {
		err := errors.Errorf("worker %q not registered", j.Name)
		failedJob(wp.Database, j, err)
		return
	}

	args := make(map[string]string)

	err := json.Unmarshal(j.Args, &args)
	if err != nil {
		err = errors.Wrapf(err, "job %d %q args unmarshal failed", j.ID(), j.Name)
		failedJob(wp.Database, j, err)
		return
	}

	j.State = running

	err = updateJob(wp.Database, j)
	if err != nil {
		return
	}

	job, err := worker.Spawn(args)
	if err != nil {
		failedJob(wp.Database, j, err)
		return
	}

	go func() {
		err = job.Run(context.Background()) // FIXME

		if err != nil {
			failedJob(wp.Database, j, err)
			return
		}

		j.State = done
		updateJob(wp.Database, j)

	}()
}

func updateJob(db primitives.Database, j Job) error {
	err := db.Update(&j)
	if err != nil {
		err := errors.Wrapf(err, "failed to update job %d %q", j.ID(), j.Name)
		log.Println(err)
	}
	return err
}

func failedJob(db primitives.Database, j Job, err error) {
	j.State = failed
	j.Error = err.Error()

	updateJob(db, j)
}
