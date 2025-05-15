package pipeline

import (
	"context"
	"log"
	"time"

	"github.com/yadhuvarshini/audio-processor/utils"
	"github.com/yadhuvarshini/audio-processor/model"

)

func StartExtractionWorker(ctx context.Context, pipe *Pipeline, workerCount int) {
	idleTimeout := 2 * time.Minute
	
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			log.Printf("🔍 ExtractionWorker %d started", workerID)

			for {
				select {
				case <-ctx.Done():
					log.Printf("🛑 ExtractionWorker %d shutting down", workerID)
					return

				case metadata := <-pipe.ExtractChan:
					log.Printf("🔍 ExtractionWorker %d: Extracting from chunkID=%s", workerID, metadata.ChunkID)

					keywords := utils.FakeExtraction(metadata.Transcript)

					log.Printf("📄 Extracted Keywords for user=%s: %v", metadata.UserID, keywords)

					result := model.FinalResult{
						ChunkID:    metadata.ChunkID,
						UserID:     metadata.UserID,
						SessionID:  metadata.SessionID,
						Timestamp:  metadata.Timestamp,
						Checksum:   metadata.Checksum,
						Transcript: metadata.Transcript,
						Keywords:   keywords,
					}

					pipe.StorageChan <- result

					case <-time.After(idleTimeout):
					log.Printf("⌛ Validation worker %d idle for 2 minutes, shutting down", id)
					return
				}
			}
		}(i)
	}
}