package jobs

import "gitlab.com/luizbranco/srs/primitives"

type Imager interface {
	primitives.Identifiable
	GetImageURL() string
	SetImageURL(string)
}

type ImageUploader interface {
	Upload(Imager)
}
