package pipeline

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/yadhuvarshini/audio-processsor/model"
)

func generateChunkID(chunk model.AudioChunk) string {
	// Generate a unique ID using data and timestamp
	hash := md5.Sum(append(chunk.Data, []byte(chunk.Timestamp.String())...))
	return hex.EncodeToString(hash[:])
}

// fakeTranscript simulates a transcription step
func fakeTranscript(data []byte) string {
	return "This is a fake transcript."
}

// generateChecksum simulates hashing the audio data
func generateChecksum(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

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
					checksum := generateChecksum(chunk.Data)
					transcript := fakeTranscript(chunk.Data)

					metadata := model.ChunkMetadata{
						ChunkID:    generateChunkID(chunk),
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
