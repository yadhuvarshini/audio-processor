package pipeline

import (
	"context"
	"fmt"

	"github.com/yadhuvarshini/audio-processsor/utils"
	"github.com/yadhuvarshini/audio-processsor/model"

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
					checksum := utils.generateChecksum(chunk.Data)
					transcript := utils.fakeTranscript(chunk.Data)

					metadata := model.ChunkMetadata{
						ChunkID:    utils.generateChunkID(chunk),
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
