package pipeline

import (
	"context"
	"fmt"

	"github.com/yadhuvarshini/audio-processor/utils"
	"github.com/yadhuvarshini/audio-processor/model"

)



// StartTransformationWorker transforms AudioChunk to ChunkMetadata

func StartTransformationWorker(ctx context.Context, pipe *Pipeline, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func(workerID int) {
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("🛑 Transformation worker %d stopped\n", workerID)
					return
				case chunk := <-pipe.ValidateChan:
					checksum := utils.GenerateChecksum(chunk.Data)
					transcript := utils.FakeTranscript(chunk.Data)

					metadata := model.ChunkMetadata{
						ChunkID:    utils.GenerateChunkID(chunk),
						UserID:     chunk.UserID,
						SessionID:  chunk.SessionID,
						Timestamp:  chunk.Timestamp,
						Checksum:   checksum,
						Transcript: transcript,
					}

					pipe.ExtractChan <- metadata
				}
			}
		}(i)
	}
}
