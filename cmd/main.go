package main

import (
	"log"
	"net/http"
	"github.com/yadhuvarshini/audio-processsor/api"
	"github.com/gorilla/mux"
)


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/upload", api.UploadHandler).Methods("POST")

	log.Println("🚀 Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
