package utils

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/yadhuvarshini/audio-processor/model"
)

func GenerateChunkID(chunk model.AudioChunk) string {
	// Generate a unique ID using data and timestamp
	hash := md5.Sum(append(chunk.Data, []byte(chunk.Timestamp.String())...))
	return hex.EncodeToString(hash[:])
}

// fakeTranscript simulates a transcription step
func FakeTranscript(data []byte) string {
	return "This is a fake transcript."
}

// generateChecksum simulates hashing the audio data
func GenerateChecksum(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}