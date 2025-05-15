package api

import (
	// "bytes"
	// "context"
	"encoding/json"
	// "fmt"
	// "io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yadhuvarshini/audio-processor/model"
	"github.com/yadhuvarshini/audio-processor/pipeline"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all for local testing
	},
}

type UploadMetadata struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Filename  string `json:"filename"`
	Timestamp string `json:"timestamp"`
}

func HandleWebSocket(pipe *pipeline.Pipeline) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
		}
		defer conn.Close()

		// 1. Read metadata
		_, metaMsg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read metadata:", err)
			return
		}

		var meta UploadMetadata
		if err := json.Unmarshal(metaMsg, &meta); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid metadata format"))
			return
		}

		conn.WriteMessage(websocket.TextMessage, []byte("✅ Metadata received"))

		// 2. Read file data (binary)
		msgType, data, err := conn.ReadMessage()
		if err != nil || msgType != websocket.BinaryMessage {
			conn.WriteMessage(websocket.TextMessage, []byte("Failed to read audio binary"))
			return
		}

		// 3. Send chunk to pipeline
		timestamp, _ := time.Parse(time.RFC3339, meta.Timestamp)
		chunk := model.AudioChunk{
			UserID:    meta.UserID,
			SessionID: meta.SessionID,
			Timestamp: timestamp,
			Data:      data,
		}

		select {
		case pipe.IngestChan <- chunk:
			conn.WriteMessage(websocket.TextMessage, []byte("📥 Chunk received and queued"))
		default:
			conn.WriteMessage(websocket.TextMessage, []byte("❌ Pipeline is full, try again later"))
			return
		}
	
	}
}
