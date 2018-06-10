package s3img

import (
	"errors"

	"gitlab.com/luizbranco/srs/primitives"
	"gitlab.com/luizbranco/srs/worker/jobs"
)

type S3img struct {
	AWSID     string
	AWSSecret string
	AWSBucket string
	AWSRegion string

	Worker   primitives.Worker
	Database primitives.Database
}

func (s3 *S3img) Upload(i jobs.Imager) error {

	return errors.New("not implemented")
}
