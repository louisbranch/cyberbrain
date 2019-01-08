package worker

import "gitlab.com/luizbranco/cyberbrain/primitives"

// ImageResizer allows an Imager to have its image resized to fit specific dimensions.
type ImageResizer interface {
	Resize(i primitives.Imager, name string, width int, height int) error
}
