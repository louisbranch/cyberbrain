package urlbuilder

import (
	"fmt"
	"strings"

	"github.com/luizbranco/srs"
	"github.com/pkg/errors"
	"github.com/speps/go-hashids"
)

type URLBuilder struct {
	hashid *hashids.HashID
}

func New() (*URLBuilder, error) {
	hd := &hashids.HashIDData{
		Alphabet:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		Salt:      "s3cr3t", // FIXME: load from env var
		MinLength: 5,
	}

	h, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize hashids")
	}

	return &URLBuilder{hashid: h}, nil
}

func (ub *URLBuilder) Path(method string, r srs.Identifiable, params ...srs.Identifiable) (string, error) {
	var qs []string

	var id string
	var err error

	if r != nil {
		id, err = ub.EncodeID(r.ID())
		if err != nil {
			return "", errors.Wrapf(err, "invalid path for record %s", r)
		}
	}

	for _, r := range params {
		// FIXME check for nil r
		slug, err := ub.EncodeID(r.ID())
		if err != nil {
			return "", errors.Wrapf(err, "invalid path for record %s", r)
		}

		qs = append(qs, fmt.Sprintf("%s=%s", r.Type(), slug))
	}

	q := strings.Join(qs, "&")

	switch method {
	case "INDEX":

		return fmt.Sprintf("/%ss/%s", r.Type(), id), nil
	case "NEW":
		return fmt.Sprintf("/%ss/new?%s", r.Type(), q), nil
	case "SHOW":
		return fmt.Sprintf("/%ss/%s", r.Type(), id), nil
	default:
		return fmt.Sprintf("/%s", r.Type()), nil
	}
}

func (ub *URLBuilder) EncodeID(id srs.ID) (string, error) {
	i := int(id)
	return ub.hashid.Encode([]int{i})
}

func (ub *URLBuilder) ParseID(hash string) (srs.ID, error) {
	ids := ub.hashid.Decode(hash)
	if len(ids) == 0 {
		return 0, errors.Errorf("invalid id for %s", hash)
	}

	id := ids[0]

	return srs.ID(id), nil
}
