package pipeline

import (
	"strings"
	"context"
	"log"

	"github.com/yadhuvarshini/audio-processor/util"

)

func StartExtractionWorker(ctx context.Context, pipe *Pipeline, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			log.Printf("🔍 ExtractionWorker %d started", workerID)

			for {
				select {
				case <-ctx.Done():
					log.Printf("🛑 ExtractionWorker %d shutting down", workerID)
					return

				case metadata := <-pipe.TransformChan:
					log.Printf("🔍 ExtractionWorker %d: Extracting from chunkID=%s", workerID, metadata.ChunkID)

					keywords := fakeExtraction(metadata.Transcript)

					log.Printf("📄 Extracted Keywords for user=%s: %v", metadata.UserID, keywords)

					// You can also send this to another channel if there's a storage/logging stage
					// pipe.StorageChan <- finalResultStruct
				}
			}
		}(i)
	}
}