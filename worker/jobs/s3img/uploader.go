package s3img

import (
	"strconv"

	"github.com/pkg/errors"
	"gitlab.com/luizbranco/srs/primitives"
	"gitlab.com/luizbranco/srs/worker/jobs"
)

const name = "s3img"

type Worker struct {
	AWSID     string
	AWSSecret string
	AWSBucket string
	AWSRegion string

	Database   primitives.Database
	WorkerPool primitives.WorkerPool
}

func (s3 *Worker) Register() error {
	err := s3.WorkerPool.Register(name, s3)
	if err != nil {
		return errors.Wrap(err, "failed to register s3 worker")
	}

	return nil
}

func (s3 *Worker) Upload(i jobs.Imager) error {
	id := strconv.Itoa(int(i.ID()))
	t := i.Type()

	args := map[string]string{
		"id":   id,
		"type": t,
	}

	err := s3.WorkerPool.Enqueue(name, args)
	if err != nil {
		return errors.Wrapf(err, "failed to enqueue s3 worker %s %s", t, id)
	}

	return nil
}

func (s3 *Worker) Spawn(args map[string]string) (primitives.Job, error) {
	return nil, errors.New("not implemented")
}
