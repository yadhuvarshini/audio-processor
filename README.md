# 🎧 Audio Processor -- Distributed Audio Stream Metadata Pipeline (Go)

### 📌 1. Assignment Objectives

This assignment required building a **concurrent audio metadata processing system** in Go with the following goals:

- Ingest audio streams (chunks) with metadata (user_id, session_id).
- Process each chunk through a multi-stage pipeline:
  - Ingestion → Validation → Transformation → Extraction → Storage.
- Store metadata persistently as JSON files.
- Provide HTTP API endpoints for:
  - Uploading chunks.
  - Querying metadata by `chunk_id`, `session_id`, and `user_id`.
  - Listing all known IDs.
- Ensure concurrency with goroutines and worker pools.
- Implement **graceful shutdown** using OS signals.
- Use Go's testing framework:
  - Table-driven unit tests.
  - Integration test.
  - Benchmarks.


## ✅ 2. Functional Requirements & Implementation

| Feature | Implemented |
|--------|-------------|
| Multipart file upload | ✅ `/upload` endpoint |
| Multi-stage pipeline with channels | ✅ All 5 stages implemented |
| Thread-safe metadata store | ✅ `sync.RWMutex`, JSON persistence |
| Metadata persistence | ✅ JSON file per chunk |
| Lookup by ID | ✅ Endpoints for chunk, session, and user IDs |
| List all IDs | ✅ `/ids` endpoint |
| Concurrency | ✅ Worker pools with goroutines |
| Graceful shutdown | ✅ Context and OS signals |
| Testing | ✅ Unit, integration, benchmark |



## 🧠 3. Architecture and Design Decisions

### 🔁 Pipeline Design

Each processing stage runs concurrently and communicates via typed Go channels:
```text
upload -> ingestionCh -> validationCh -> transformCh -> extractionCh -> storageCh

```

Each stage performs a specific transformation or validation, promoting separation of concerns.

### 🧵 Concurrency

Worker pools were implemented for each stage using `go func()` and channels. Context cancellation ensures all goroutines exit gracefully on shutdown.

### 📁 Storage

A custom `MetadataStore` was implemented:

-   In-memory map for fast access.

-   `sync.RWMutex` for thread safety.

-   JSON files for durable storage.

* * * * *

🧪 4. Testing Approach
----------------------

### ✅ Unit Testing

-   Table-driven unit tests for:

    -   `GenerateChecksum`

    -   `GenerateChunkID`

    -   `FakeTranscript`

### 🔄 Integration Test

-   End-to-end test simulating:

    -   File upload via HTTP

    -   Pipeline processing

    -   Metadata persistence

### ⚡ Benchmarking

-   Benchmark test for `GenerateChecksum` using Go's `testing.B`.

* * * * *

📁 5. Folder and File Details
-----------------------------

```
audio-processor/
├── api/
│   └── api.go              # HTTP handlers for upload and query APIs
│
├── model/
│   └── model.go            # Structs for AudioChunk and FinalResult
│
├── pipeline/
│   ├── pipeline.go         # Implementation of the 5 pipeline stages
│   └── pipeline_test.go    # Unit + benchmark tests
│
├── storage/
│   └── store.go            # MetadataStore: map + mutex + JSON files
│
├── utils/
│   └── utils.go            # Helper functions: checksum, transcript, ID gen
│
├── main.go                 # App entry point, router, server, shutdown
├── go.mod / go.sum         # Module management

```

* * * * *

🔍 6. Known Limitations & Future Work
-------------------------------------

### ❗ Limitations

-   Transcript generation is simulated (not real NLP/ML).

-   No authentication or authorization on APIs.

-   JSON file storage is suitable for small scale only.

### 🛠️ Future Enhancements

-   Replace fake transcript with real speech-to-text.

-   Add user auth (JWT/session tokens).

-   Use database (PostgreSQL/BadgerDB) instead of flat files.

-   Pagination/filtering on list endpoints.

-   Add Prometheus/Grafana metrics for observability.

* * * * *

🔑 7. How to Run the Project
----------------------------

### ✅ Prerequisites

-   Go 1.19+

-   Git

### 🛠️ Setup

```
git clone https://github.com/yadhuvarshini/audio-processor.git
cd audio-processor
go mod tidy

```

### 🚀 Run the Server

```
go run main.go

```

Server will start on port `8080`.

* * * * *

🌐 8. API Endpoints
-------------------

### 🟢 POST `/upload`

-   Accepts audio file (multipart) with:

    -   `user_id` (form field)

    -   `session_id` (form field)

```
curl -F "file=@test.wav" -F "user_id=abc" -F "session_id=s1" http://localhost:8080/upload

```

* * * * *

### 🔵 GET `/chunks/{chunk_id}`

Get metadata for a specific chunk.

```
curl http://localhost:8080/chunks/abc123

```

* * * * *

### 🔵 GET `/sessions/{user_id}`

Get all sessions for a given user.

```
curl http://localhost:8080/sessions/user123

```

* * * * *

### 🔵 GET `/ids`

List all chunk IDs, session IDs, and user IDs.

```
curl http://localhost:8080/ids

```

* * * * *

🧩 9. Important Code Highlights
-------------------------------

### 🔐 GenerateChecksum

```
func GenerateChecksum(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

```

### 🧠 Fake Transcript

```
func FakeTranscript(data []byte) string {
	return "This is a fake transcript."
}

```

### 🧪 Unit Test for Checksum

```
func TestGenerateChecksum(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"text1", []byte("hello"), "5d41402abc4b2a76b9719d911017c592"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateChecksum(tt.input)
			if got != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, got)
			}
		})
	}
}

```

### 🧪 Integration Test: Upload Handler

```
func TestUploadHandler_Integration(t *testing.T) {
	req := httptest.NewRequest("POST", "/upload", formData)
	// Set headers and body...

	// Run pipeline workers concurrently
	go api.UploadHandler(rr, req, pipeline)

	time.Sleep(2 * time.Second)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", rr.Code)
	}
}

```

### 🧵 Graceful Shutdown

```
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt)
<-c
cancel()
server.Shutdown(ctx)

```

* * * * *

🏁 10. Conclusion
-----------------

This project demonstrates a fully functional, concurrent Go application that processes audio stream metadata using channels, goroutines, and file-based persistence. It was completed in ~12 hours by a single developer and covers concurrency, pipeline design, HTTP APIs, and testing comprehensively.

* * * * *

*Developed with ❤️ by [Yadhuvarshini](https://github.com/yadhuvarshini/)*

```
