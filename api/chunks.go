package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yadhuvarshini/audio-processor/store"
)

var metadataStore *store.MetadataStore

// Inject store (do this in main)
func SetStore(s *store.MetadataStore) {
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
