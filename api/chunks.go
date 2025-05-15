package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yadhuvarshini/audio-processor/storage"
)

var metadataStore *storage.MetadataStore

// Inject store (do this in main)
func SetStore(s *storage.MetadataStore) {
	metadataStore = s
}

func GetChunkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chunkID := vars["id"]

	result, ok := metadataStore.GetByChunkID(chunkID)
	if !ok {
		http.Error(w, "Chunk not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func GetUserChunksHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	results := metadataStore.GetByUserID(userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

