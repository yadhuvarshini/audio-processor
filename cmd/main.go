package main

import (
	"log"
	"net/http"
	"context"

	"github.com/gorilla/mux"

	"github.com/yadhuvarshini/audio-processsor/api"
	"github.com/yadhuvarshini/audio-processsor/pipeline"
)


func main() {
	r := mux.NewRouter()
	ctx := context.Background()

	r.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		pipe := pipeline.NewPipeline(ctx) // Create or initialize the pipeline instance
		pipeline.StartIngestionWorker(ctx, pipe, 5)
		pipeline.StartValidationWorkers(ctx, pipe, 5)    // 5 validation workers
		pipeline.StartTransformationWorker(ctx, pipe, 3)

		api.UploadHandler(w, r, pipe)
	}).Methods("POST")

	log.Println("🚀 Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

