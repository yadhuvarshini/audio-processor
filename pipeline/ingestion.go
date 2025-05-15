package pipeline

import (
	"context"
	"log"
	"time"

	"github.com/yadhuvarshini/audio-processsor/model"
)

// Starting the Worker pool
// Receives chunk from /upload API (you already implemented this!)
// Pushes it into IngestChan
// IngestChan -> ValidateChan -> TransformChan -> ExtractChan -> StorageChan
// Each channel has a buffer of 100

// Receives from IngestChan
// Simulates delay (e.g., network lag)
// Passes to ValidateChan

//In real systems, ingestion might:
// Do basic formatting or cleanup
// Enrich with session info
// Rate-limit users
// Add chunk IDs or timestamps

func StartIngestionWorker(ctx context.Context, p *Pipeline, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go func(id int) {
			log.Printf("Ingestion worker %d started\n", id)
			for {
				select {
				case <-ctx.Done():
					log.Printf("Ingestion worker %d stopping\n", id)
					return

				case chunk := <-p.IngestChan:
					// Simulate some processing delay
					time.Sleep(50 * time.Millisecond) // 50ms = 0.05s
					log.Printf("Ingestion worker %d processing chunk from user=%s session=%s\n", id, chunk.UserID, chunk.SessionID)
					// Pass the chunk to the next stage
					p.ValidateChan <- chunk
				}
			}
		}(i)
	}
}

//COPILOT NOTE:
// In a real-world scenario, the ingestion worker might also:
// Validate the audio format
// Check for duplicate chunks
// Enrich metadata (e.g., adding a unique chunk ID)
// For example, if the audio is in a different format (e.g., WAV instead of MP3), the ingestion worker might convert it to the expected format before passing it to the next stage.
// This is a simplified example, and in a real-world application, you would likely have more complex logic for handling different audio formats, error handling, and other considerations.