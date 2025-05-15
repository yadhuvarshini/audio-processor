package pipeline

import (
	"context"
	"log"

	"github.com/yadhuvarshini/audio-processor/model"
	"github.com/yadhuvarshini/audio-processor/utils"
)


// StartTransformationWorker reads from ValidateChan, transforms chunks, and sends to TransformChan
func StartTransformationWorker(ctx context.Context, pipe *Pipeline, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			log.Printf("⚙️  TransformationWorker %d started", workerID)
			log.Println("DEBUG : Tr", pipe.ValidateChan)

			for {
				select {
				case <-ctx.Done():
					log.Printf("🛑 TransformationWorker %d shutting down", workerID)
					return

				case chunk := <-pipe.ValidateChan:
					log.Printf("🧪 TransformationWorker %d: Processing chunk from user=%s", workerID, chunk.UserID)

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

					log.Printf("✅ TransformationWorker %d: Sent metadata for user=%s", workerID, chunk.UserID)
				}
			}
		}(i)
	}
}
