package jobs

import "gitlab.com/luizbranco/cyberbrain/primitives"

type Imager interface {
	primitives.Identifiable
	GetImageURL() string
	SetImageURL(string)
}

type ImageUploader interface {
	Upload(Imager, string) error
}
