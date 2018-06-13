package s3img

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"gitlab.com/luizbranco/srs/db"
	"gitlab.com/luizbranco/srs/primitives"
)

type Job struct {
	ID   primitives.ID
	Type string

	awsBucket string
	db        primitives.Database
	client    *s3.S3
}

func (j *Job) Run(ctx context.Context) error {
	if j.Type != "card" {
		return errors.Errorf("%q upload type not implemented", j.Type)
	}

	card, err := db.FindCard(j.db, j.ID)
	if err != nil {
		return errors.Wrapf(err, "failed to find record %q %d", j.Type, j.ID)
	}

	res, err := http.Get(card.GetImageURL())
	if err != nil {
		return errors.Wrapf(err, "failed to read image %q %d", j.Type, j.ID)
	}

	defer res.Body.Close()

	buff := bytes.NewBuffer(nil)

	_, err = io.Copy(buff, res.Body)
	if err != nil {
		return errors.Wrapf(err, "failed to copy image %q %d", j.Type, j.ID)
	}

	b := buff.Bytes()

	var name string

	mime := http.DetectContentType(b)
	switch mime {
	case "image/jpeg":
		name = fmt.Sprintf("%ss/%d.jpg", j.Type, j.ID)
	default:
		return errors.Wrapf(err, "mime-type %q cannot be used for image %q %d", mime, j.Type, j.ID)
	}

	r := bytes.NewReader(b)

	_, err = j.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(j.awsBucket),
		Key:    aws.String(name),
		Body:   r,
	})

	if err != nil {
		return errors.Wrapf(err, "failed to upload image to AWS S3 %q %d", j.Type, j.ID)
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", j.awsBucket, name)

	card.SetImageURL(url)

	err = j.db.Update(card)
	if err != nil {
		return errors.Wrapf(err, "failed to update image url %q %d %q", j.Type, j.ID, url)
	}

	return nil
}
