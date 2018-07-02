package offline

import (
	"log"

	"gitlab.com/luizbranco/cyberbrain/worker"
)

type ImageOfflineResizer struct{}

func (w *ImageOfflineResizer) Resize(i worker.Imager, name string, width int, height int) error {
	log.Printf("image resize called for %s\n", name)
	return nil
}
