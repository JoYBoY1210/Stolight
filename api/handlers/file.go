package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joyboy1210/stolight/models"
)

func ListFilesInBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Query().Get("bucket")
	if bucketName == "" {
		http.Error(w, "Bucket name is required to list all the files in it", http.StatusBadRequest)
		return
	}
	bucket, err := models.GetBucketByName(bucketName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get bucket: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	files, err := models.GetFilesByBucketID(bucket.ID)
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status": "success",
		"files":  files,
	}
	json.NewEncoder(w).Encode(response)
}

