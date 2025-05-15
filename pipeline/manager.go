package pipeline

import (

	"context"  //concurrency control and exit
	"github.com/yadhuvarshini/audio-processsor/model"
	
	)
type Pipeline struct {
	IngestChan chan model.AudioChunk
	ValidateChan chan model.AudioChunk
	TransformChan chan model.AudioChunk
	ExtractChan chan model.ChunkMetadata
	StorageChan chan model.ChunkMetadata
}

func NewPipeline(ctx context.Context) *Pipeline {
	return &Pipeline{
		IngestChan:   make(chan model.AudioChunk, 100),
		ValidateChan: make(chan model.AudioChunk, 100),
		TransformChan: make(chan model.AudioChunk, 100),
		ExtractChan:  make(chan model.ChunkMetadata, 100),
		StorageChan:  make(chan model.ChunkMetadata, 100),
	}
}



