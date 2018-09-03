package urlbuilder

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/speps/go-hashids"
	"gitlab.com/luizbranco/cyberbrain/primitives"
)

type URLBuilder struct {
	hashid *hashids.HashID
}

type Slugger interface {
	Slug() string
}

func New(salt string) (*URLBuilder, error) {
	hd := &hashids.HashIDData{
		Alphabet:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		Salt:      salt,
		MinLength: 5,
	}

	h, err := hashids.NewWithData(hd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize hashids")
	}

	return &URLBuilder{hashid: h}, nil
}

func (ub *URLBuilder) Path(method string, idn primitives.Identifiable,
	params ...primitives.Identifiable) (string, error) {

	var qs []string

	var id string
	var err error

	if idn != nil {
		id, err = ub.EncodeID(idn.ID())
		if err != nil {
			return "", errors.Wrapf(err, "invalid path for record %s", idn)
		}
	}

	var prefix string

	for _, r := range params {
		if r == nil {
			continue
		}

		slug, err := ub.EncodeID(r.ID())
		if err != nil {
			return "", errors.Wrapf(err, "invalid path for record %s", r)
		}

		if r.Type() == "deck" {
			prefix = "/decks/" + slug
		} else {
			qs = append(qs, fmt.Sprintf("%s=%s", r.Type(), slug))
		}
	}

	var q string

	if len(qs) > 0 {
		q = "?" + strings.Join(qs, "&")
	}

	slug := idn.Type()
	slugger, ok := idn.(Slugger)
	if ok {
		slug = slugger.Slug()
	}

	switch method {
	case "INDEX":
		return fmt.Sprintf("%s/%ss", prefix, slug), nil
	case "NEW":
		return fmt.Sprintf("%s/%ss/new%s", prefix, slug, q), nil
	case "SHOW":
		return fmt.Sprintf("%s/%ss/%s", prefix, slug, id), nil
	case "CREATE":
		return fmt.Sprintf("%s/%ss", prefix, slug), nil
	default:
		return fmt.Sprintf("%s/%s", prefix, slug), nil
	}
}

func (ub *URLBuilder) EncodeID(id primitives.ID) (string, error) {
	i := int(id)
	return ub.hashid.Encode([]int{i})
}

func (ub *URLBuilder) ParseID(hash string) (primitives.ID, error) {
	ids, err := ub.hashid.DecodeWithError(hash)
	if err != nil || len(ids) == 0 {
		return 0, errors.Wrapf(err, "invalid id for %s", hash)
	}

	id := ids[0]

	return primitives.ID(id), nil
}
