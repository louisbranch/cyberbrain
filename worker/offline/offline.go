package offline

import (
	"log"

	"gitlab.com/luizbranco/cyberbrain/primitives"
)

// ImageOfflineResizer implements ImageResizer interface.
type ImageOfflineResizer struct{}

// Resize is a no-op function for offline work.
func (w *ImageOfflineResizer) Resize(i primitives.Imager, name string, width int, height int) error {
	log.Printf("image resize called for %s\n", name)
	return nil
}
