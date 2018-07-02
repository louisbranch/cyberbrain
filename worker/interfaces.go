package worker

import "gitlab.com/luizbranco/cyberbrain/primitives"

type Imager interface {
	primitives.Identifiable
	GetImageURL() string
	SetImageURL(string)
}

type ImageResizer interface {
	Resize(i Imager, name string, width int, height int) error
}

type CardScheduler interface {
	Schedule(deckID, cardID primitives.ID, correct bool) error
}
