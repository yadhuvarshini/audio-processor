package pipeline

import (
	"context"
	"log"
	"time"
	
)

// Starts validation workers
func StartValidationWorkers(ctx context.Context, p *Pipeline, workerCount int) {
	idleTimeout := 2 * time.Minute

	for i := 0; i < workerCount; i++ {
		go func(id int) {
			log.Printf("✅ Validation worker %d started", id)
			for {
				select {
				case <-ctx.Done():
					log.Printf("🛑 Validation worker %d stopped", id)
					return
				case chunk := <-p.ValidateChan:
					// Simulate validation delay
					time.Sleep(30 * time.Millisecond)

					if chunk.UserID == "" || chunk.SessionID == "" || chunk.Timestamp.IsZero() {
						log.Printf("❌ Validation failed: missing metadata [worker=%d]", id)
						continue // drop or log invalid chunks
					}

					log.Printf("✅ Chunk validated [worker=%d] user=%s", id, chunk.UserID)

					// Send to next stage
					p.TransformChan <- chunk

				case <-time.After(idleTimeout):
					log.Println("⌛ Ingestion worker %d idle for 2 minutes, shutting down", id)
					return
				}
			}
		}(i)
	}
}
