package model

import (
	"time" //for timestamp
)

// Raw audio + metadata sent by the client to the server
type AudioChunk struct {
	UserID    string
	SessionID string
	Timestamp time.Time
	Data      []byte
}

//processes metadata we extract

type ChunkMetadata struct {
	ChunkID    string
	UserID     string
	SessionID  string
	Timestamp  time.Time
	Checksum   string // Checksum of the audio data (fake)
	Transcript string // Transcription of the audio data (fake)
}
