package pipeline

import (
	"context"
	"log"
	"time"
	
)

// Starts validation workers
func StartValidationWorkers(ctx context.Context, p *Pipeline, workerCount int) {	

	for i := 0; i < workerCount; i++ {
		go func(id int) {
			log.Printf("✅ Validation worker %d started", id)
			for {
				log.Printf("DEBUG: ingestion - 1")
				select {
				case <-ctx.Done():
					log.Printf("🛑 Validation worker %d stopped", id)
					return
				case chunk := <-p.ValidateChan:
					// Simulate validation delay
					log.Println("Debug - 4")
					time.Sleep(30 * time.Millisecond)

					if chunk.UserID == "" || chunk.SessionID == "" || chunk.Timestamp.IsZero() {
						log.Printf("❌ Validation failed: missing metadata [worker=%d]", id)
						continue // drop or log invalid chunks
					}

					log.Printf("✅ Chunk validated [worker=%d] user=%s", id, chunk.UserID)

					// Send to next stage
					p.TransformChan <- chunk

				
				}
			}
		}(i)
	}
}
