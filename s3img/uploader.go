package s3img

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"gitlab.com/luizbranco/cyberbrain/primitives"
)

// S3Img implements the ImageUploader interface.
type S3Img struct {
	awsBucket string
	client    *s3.S3
}

// New created a new S3Img uploader
func New(id, secret, region, bucket string) (*S3Img, error) {
	sess := session.Must(
		session.NewSession(&aws.Config{
			Region: aws.String(region),
		}))

	creds := credentials.NewStaticCredentials(id, secret, "")

	s3img := S3Img{
		awsBucket: bucket,
		client:    s3.New(sess, &aws.Config{Credentials: creds}),
	}

	return &s3img, nil
}

// Upload an image to S3 and update the record on the DB.
func (s3img *S3Img) Upload(ctx context.Context, i primitives.Imager, name string, img []byte) (string, error) {
	if i.Type() != "card" {
		return "", errors.Errorf("%q upload type not implemented", i.Type())
	}

	var ext string

	mime := http.DetectContentType(img)
	switch mime {
	case "image/png":
		ext = "png"
	case "image/jpeg":
		ext = "jpg"
	default:
		return "", errors.Errorf("mime-type %q cannot be used for image %q %d", mime, i.Type(), i.ID())
	}

	key := fmt.Sprintf("%ss/%s.%s", i.Type(), name, ext)

	r := bytes.NewReader(img)

	_, err := s3img.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s3img.awsBucket),
		Key:    aws.String(key),
		Body:   r,
	})

	if err != nil {
		return "", errors.Wrapf(err, "failed to upload image to AWS S3 %q %d", i.Type(), i.ID())
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s3img.awsBucket, key)

	return url, nil
}
