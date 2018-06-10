package jobs

type Imager interface {
	GetImageURL() string
	SetImageURL(string)
}

type ImageUploader interface {
	Upload(Imager)
}
