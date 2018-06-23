package resizer

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/luizbranco/cyberbrain/primitives"
	"gitlab.com/luizbranco/cyberbrain/worker"
)

const workerName = "img-resize"

type Worker struct {
	BlitlineID  string
	AWSBucket   string
	CallbackURL string
	Poll        bool

	WorkerPool primitives.WorkerPool
}

func (w *Worker) Register() error {
	if w.WorkerPool == nil {
		return errors.New("invalid worker pool")
	}

	if w.BlitlineID == "" {
		return errors.New("BlitlineID cannot be empty")
	}

	if w.AWSBucket == "" {
		return errors.New("AWSBucket cannot be empty")
	}

	err := w.WorkerPool.Register(workerName, w)
	if err != nil {
		return errors.Wrap(err, "failed to register image resize worker")
	}

	return nil
}

func (w *Worker) Resize(i worker.Imager, name string, width, height int) error {

	if w.WorkerPool == nil {
		return errors.New("invalid worker pool")
	}

	callback := fmt.Sprintf("%s/%ss/%s", w.CallbackURL, i.Type(), name)
	s3Path := fmt.Sprintf("%ss/%s.png", i.Type(), name)

	url := i.GetImageURL()

	args := JobArgs{
		ImagerID:    name,
		ImageURL:    url,
		CallbackURL: callback,
		Poll:        w.Poll,
		S3Path:      s3Path,
		Width:       width,
		Height:      height,
	}

	err := w.WorkerPool.Enqueue(workerName, args)
	if err != nil {
		return errors.Wrapf(err, "failed to enqueue image resize worker %q %q", url, name)
	}

	return nil
}

func (w *Worker) Spawn(b []byte) (primitives.Job, error) {
	args := JobArgs{}

	err := json.Unmarshal(b, &args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal args")
	}

	j := &Job{
		args:           args,
		awsBucket:      w.AWSBucket,
		blitlineAPIURL: blitlineAPIURL,
		blitlineID:     w.BlitlineID,
	}

	return j, nil
}
