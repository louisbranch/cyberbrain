package s3img

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"gitlab.com/luizbranco/srs/primitives"
	"gitlab.com/luizbranco/srs/worker/jobs"
)

const workerName = "s3img"

type Worker struct {
	AWSID     string
	AWSSecret string
	AWSBucket string
	AWSRegion string

	Database   primitives.Database
	WorkerPool primitives.WorkerPool

	client *s3.S3
}

func (w *Worker) Register() error {
	err := w.WorkerPool.Register(workerName, w)
	if err != nil {
		return errors.Wrap(err, "failed to register s3 worker")
	}

	sess := session.Must(
		session.NewSession(&aws.Config{
			Region: aws.String(w.AWSRegion),
		}))

	creds := credentials.NewStaticCredentials(w.AWSID, w.AWSSecret, "")

	w.client = s3.New(sess, &aws.Config{Credentials: creds})

	return nil
}

func (w *Worker) Upload(i jobs.Imager) error {
	id := strconv.Itoa(int(i.ID()))
	t := i.Type()

	args := map[string]string{
		"id":   id,
		"type": t,
	}

	err := w.WorkerPool.Enqueue(workerName, args)
	if err != nil {
		return errors.Wrapf(err, "failed to enqueue s3 worker %s %s", t, id)
	}

	return nil
}

func (w *Worker) Spawn(args map[string]string) (primitives.Job, error) {
	j := &Job{
		Type:      args["type"],
		client:    w.client,
		db:        w.Database,
		awsBucket: w.AWSBucket,
	}

	id, err := strconv.Atoi(args["id"])
	if err != nil {
		return nil, errors.Wrapf(err, "invalid id %q %q", j.Type, id)
	}

	j.ID = primitives.ID(id)

	return j, nil
}
