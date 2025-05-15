package main

import (
	"context"
	"log"
	"net/http"
	"encoding/json"


	"github.com/gorilla/mux"

	"github.com/yadhuvarshini/audio-processor/api"
	"github.com/yadhuvarshini/audio-processor/pipeline"
	"github.com/yadhuvarshini/audio-processor/storage"
)

func main() {
	r := mux.NewRouter()
	ctx := context.Background()
	store := storage.NewMetadataStore("data/chunks") // Folder to save JSON files

	pipe := pipeline.NewPipeline(ctx) // Create or initialize the pipeline instance
	pipeline.StartIngestionWorker(ctx, pipe, 5)
	pipeline.StartValidationWorkers(ctx, pipe, 5) // 5 validation workers
	pipeline.StartTransformationWorker(ctx, pipe, 5)
	pipeline.StartExtractionWorker(ctx, pipe, 5)
	pipeline.StartStorageWorker(ctx, pipe, store, 5) // 👈 Add this line


	r.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		api.UploadHandler(w, r, pipe)
	}).Methods("POST")

	r.HandleFunc("/chunks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		meta, ok := store.GetByChunkID(id)
		if !ok {
			http.Error(w, "Chunk not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(meta)
	}).Methods("GET")

	r.HandleFunc("/sessions/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		userID := mux.Vars(r)["user_id"]
		results := store.GetByUserID(userID)
		json.NewEncoder(w).Encode(results)
	}).Methods("GET")

	r.HandleFunc("/session-chunks/{session_id}", func(w http.ResponseWriter, r *http.Request) {
		sessionID := mux.Vars(r)["session_id"]
		results := store.GetBySessionID(sessionID)
		json.NewEncoder(w).Encode(results)
	}).Methods("GET")

	r.HandleFunc("/ids", func(w http.ResponseWriter, r *http.Request) {
	chunks, sessions, users := store.ListAllIDs()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"chunk_ids":   chunks,
		"session_ids": sessions,
		"user_ids":    users,
	})
}).Methods("GET")

	log.Println("🚀 Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
