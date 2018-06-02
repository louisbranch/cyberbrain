package urlbuilder

import (
	"fmt"

	"github.com/luizbranco/srs/web"
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

func (ub *URLBuilder) Path(method string, r web.Record, params ...web.Record) (string, error) {
	id := int(r.ID())

	slug, err := ub.hashid.Encode([]int{id})
	if err != nil {
		return "", errors.Wrapf(err, "invalid path for record %s", r)
	}

	path := fmt.Sprintf("/%ss/%s", r.Type(), slug)
	return path, nil
}

func (ub *URLBuilder) ID(hash string) (web.ID, error) {
	ids := ub.hashid.Decode(hash)
	if len(ids) == 0 {
		return 0, errors.Errorf("invalid id for %s", hash)
	}

	id := ids[0]

	return web.ID(id), nil
}
