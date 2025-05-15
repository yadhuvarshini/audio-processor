package pipeline

import (
	"context"
	"log"
	"time"

	"github.com/yadhuvarshini/audio-processor/model"
	"github.com/yadhuvarshini/audio-processor/utils"
)


func StartTransformationWorker(ctx context.Context, pipe *Pipeline, workerCount int) {
	idleTimeout := 2 * time.Minute

	for i := 0; i < workerCount; i++ {
		go func(id int) {
			log.Printf("⚙️  TransformationWorker %d started", id)
			log.Println("DEBUG : Tr", pipe.TransformChan)

			for {
				select {
				case <-ctx.Done():
					log.Printf("🛑 TransformationWorker %d shutting down", id)
					return

				case chunk := <-pipe.TransformChan:
					log.Printf("🧪 TransformationWorker %d: Processing chunk from user=%s", id, chunk.UserID)

					// Simulate transformation logic
					checksum := utils.GenerateChecksum(chunk.Data)
					transcript := utils.FakeTranscript(chunk.Data)
					chunkID := utils.GenerateChunkID(chunk)

					metadata := model.ChunkMetadata{
						ChunkID:    chunkID,
						UserID:     chunk.UserID,
						SessionID:  chunk.SessionID,
						Timestamp:  chunk.Timestamp,
						Checksum:   checksum,
						Transcript: transcript,
					}

					// Send transformed metadata to next stage
					pipe.ExtractChan <- metadata

					log.Printf("✅ TransformationWorker %d: Sent metadata for user=%s", id, chunk.UserID)

					case <-time.After(idleTimeout):
					log.Printf("⌛ Validation worker %d idle for 2 minutes, shutting down", id)
					return
				}
			}
		}(i)
	}
}
