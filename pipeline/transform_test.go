package pipeline

import (
	"testing"
	"time"

	"github.com/yadhuvarshini/audio-processor/utils"
	"github.com/yadhuvarshini/audio-processor/model"
)

func TestGenerateChecksum(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"simple text", []byte("hello"), "5d41402abc4b2a76b9719d911017c592"},
		{"another text", []byte("world"), "7d793037a0760186574b0282f2f435e7"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := utils.GenerateChecksum(tt.input)
			if actual != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, actual)
			}
		})
	}
}

func TestGenerateChunkID_Unique(t *testing.T) {
	chunk1 := model.AudioChunk{
		Data:      []byte("audio1"),
		Timestamp: time.Now(),
	}
	chunk2 := model.AudioChunk{
		Data:      []byte("audio2"),
		Timestamp: time.Now().Add(time.Second),
	}

	id1 := utils.GenerateChunkID(chunk1)
	id2 := utils.GenerateChunkID(chunk2)

	if id1 == id2 {
		t.Errorf("Chunk IDs should be unique but got the same: %s", id1)
	}
}

func TestFakeTranscript(t *testing.T) {
	data := []byte("test audio")
	expected := "This is a fake transcript."
	actual := utils.FakeTranscript(data)

	if actual != expected {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}

func BenchmarkGenerateChecksum(b *testing.B) {
	data := []byte("some sample audio data")
	for i := 0; i < b.N; i++ {
		_ = utils.GenerateChecksum(data)
	}
}
