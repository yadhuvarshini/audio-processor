package api

import (
	"log"
	"net/http"
	"time"
	"io"

	"github.com/yadhuvarshini/audio-processsor/utils"
	"github.com/yadhuvarshini/audio-processsor/model"
	"github.com/yadhuvarshini/audio-processsor/pipeline"
)


func UploadHandler(w http.ResponseWriter, r *http.Request, pipe *pipeline.Pipeline) {
	

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	} // Max 10MB

	file, _, err := r.FormFile("file") //getting a file
	if err != nil {  //if it has error
		utils.RespondWithError(w, http.StatusBadRequest, "Missing File")
		return
	}
	defer file.Close() //bad file has been closed

	//metadata is procured
	userID := r.FormValue("user_id")
	sessionID := r.FormValue("session_id")
	timestamp := time.Now()
	
	log.Printf("DEBUG: user_id=%q, session_id=%q\n", userID, sessionID)


	if userID == "" || sessionID == ""  {
        utils.RespondWithError(w, http.StatusBadRequest, "Missing metadata fields")
        return
    }

	data, err := io.ReadAll(file)  //reading the file to pass it to file in audiochunk data structure
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to read file")
		return
	}

	chunk := &model.AudioChunk{
		UserID:    userID,
		SessionID: sessionID,
		Timestamp: timestamp,
		Data:      data,
		 
	}

	pipe.IngestChan <- *chunk

	log.Printf("📥 Received chunk from user=%s session=%s\n", userID, sessionID)

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "success", "message": "File uploaded successfully"})
}