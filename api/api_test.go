package api

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yadhuvarshini/audio-processor/pipeline"
	"github.com/yadhuvarshini/audio-processor/storage"
)

func TestUploadHandler_Integration(t *testing.T) {
	// Prepare fake multipart request
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.wav")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("fake audio")) // simulate audio data
	writer.WriteField("user_id", "test-user")
	writer.WriteField("session_id", "test-session")
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()

	// Setup pipeline and store
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := storage.NewMetadataStore("testdata") // use temp dir for tests
	pipe := pipeline.NewPipeline(ctx)

	// Start minimal workers
	pipeline.StartIngestionWorker(ctx, pipe, 1)
	pipeline.StartValidationWorkers(ctx, pipe, 1)
	pipeline.StartTransformationWorker(ctx, pipe, 1)
	pipeline.StartExtractionWorker(ctx, pipe, 1)
	pipeline.StartStorageWorker(ctx, pipe, store, 1)

	// Invoke the handler
	go UploadHandler(rr, req, pipe)

	// Wait for processing (tunable)
	time.Sleep(2 * time.Second)

	// Check response
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", rr.Code)
	}

	// Optionally check: has the store saved any chunks?
	// Example: fetch all chunk IDs or validate the saved file
	found := false
	store.ListAllIDs()
	chunks, _, _ := store.ListAllIDs()
	for _, cid := range chunks {
		meta, ok := store.GetByChunkID(cid)
		if ok && meta.UserID == "test-user" && meta.SessionID == "test-session" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("Expected metadata to be stored but it was not found")
	}
}
