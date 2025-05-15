package pipeline

import (
	"context"

	"github.com/yadhuvarshini/audio-processor/storage"
)

func StartStorageWorker(ctx context.Context, pipe *Pipeline, store *storage.MetadataStore, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case result := <-pipe.StorageChan:
					err := store.Save(result)
					if err != nil {
						// Optionally log error
					}
				}
			}
		}()
	}
}
