package api

import (
	"log"
	"net/http"
	"github.com/yadhuvarshini/audio-processsor/utils"
)


func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	} // Max 10MB

	file, _, err := r.FormFile("file")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing File")
		return
	}
	defer file.Close()

	userID := r.FormValue("user_id")
	sessionID := r.FormValue("session_id")
	timestamp := r.FormValue("timestamp")

	if userID == "" || sessionID == "" || timestamp == "" {
        utils.RespondWithError(w, http.StatusBadRequest, "Missing metadata fields")
        return
    }

	// _, err := io.ReadAll(file)
    // if err != nil {
    //     utils.RespondWithError(w, http.StatusInternalServerError, "Failed to read file")
    //     return
    // }

	log.Printf("📥 Received chunk from user=%s session=%s at %s\n", userID, sessionID, timestamp)

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "success", "message": "File uploaded successfully"})
}