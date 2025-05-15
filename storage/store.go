package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/yadhuvarshini/audio-processor/model"
)

type MetadataStore struct {
	mu       sync.RWMutex
	chunks   map[string]model.FinalResult // chunk_id -> metadata
	basePath string                       // for storing JSON files
}

func NewMetadataStore(basePath string) *MetadataStore {
	os.MkdirAll(basePath, os.ModePerm) // ensure dir exists
	return &MetadataStore{
		chunks:   make(map[string]model.FinalResult),
		basePath: basePath,
	}
}

func (s *MetadataStore) Save(result model.FinalResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Save in memory
	s.chunks[result.ChunkID] = result

	// Save to disk
	filePath := filepath.Join(s.basePath, result.ChunkID + ".json")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(result)
}

func (s *MetadataStore) GetByChunkID(chunkID string) (model.FinalResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result, ok := s.chunks[chunkID]
	return result, ok
}

func (s *MetadataStore) GetByUserID(userID string) []model.FinalResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var results []model.FinalResult
	for _, v := range s.chunks {
		if v.UserID == userID {
			results = append(results, v)
		}
	}
	return results
}

func (s *MetadataStore) GetBySessionID(sessionID string) []model.FinalResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var results []model.FinalResult
	for _, v := range s.chunks {
		if v.SessionID == sessionID {
			results = append(results, v)
		}
	}
	return results
}

