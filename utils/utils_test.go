package utils_test

import (
	"testing"
	"time"

	"github.com/yadhuvarshini/audio-processor/model"
	"github.com/yadhuvarshini/audio-processor/utils"
)

func TestGenerateChunkID(t *testing.T) {
	tests := []struct {
		name   string
		chunk  model.AudioChunk
		unique bool
	}{
		{"UniqueIDs", model.AudioChunk{Data: []byte("A"), Timestamp: time.Now()}, true},
		{"UniqueIDsDifferentTime", model.AudioChunk{Data: []byte("A"), Timestamp: time.Now().Add(time.Second)}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id1 := utils.GenerateChunkID(tt.chunk)
			id2 := utils.GenerateChunkID(tt.chunk)
			if tt.unique && id1 != id2 {
				t.Errorf("Expected same ID, got different: %s vs %s", id1, id2)
			}
		})
	}
}
