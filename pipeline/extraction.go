package pipeline

import (
	"context"
	"log"

	"github.com/yadhuvarshini/audio-processor/utils"
	"github.com/yadhuvarshini/audio-processor/model"

)

func StartExtractionWorker(ctx context.Context, pipe *Pipeline, workerCount int) {
	
	for i := 0; i < workerCount; i++ {
		go func(id int) {
			log.Printf("🔍 ExtractionWorker %d started", id)

			for {
				select {
				case <-ctx.Done():
					log.Printf("🛑 ExtractionWorker %d shutting down", id)
					return

				case metadata := <-pipe.ExtractChan:
					log.Printf("🔍 ExtractionWorker %d: Extracting from chunkID=%s", id, metadata.ChunkID)

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

				}
			}
		}(i)
	}
}