package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joyboy1210/stolight/models"
)

func ListFilesInBucketHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bucketName := r.PathValue("bucket")
	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}
	bucket, err := models.GetBucketByName(bucketName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get bucket: %s", err.Error()), http.StatusNotFound)
		return
	}

	files, err := models.GetFilesByBucketID(bucket.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get files in bucket: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status": "success",
		"files":  files,
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {

	bucketName := r.PathValue("bucket")
	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}
	fileId := r.PathValue("fileId")
	if fileId == "" {
		http.Error(w, "File ID is required", http.StatusBadRequest)
		return
	}

	bucketMeta, err := models.GetBucketByName(bucketName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get bucket: %s", err.Error()), http.StatusNotFound)
		return
	}

	fileMeta, err := models.GetFileByID(fileId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get file metadata: %s", err.Error()), http.StatusNotFound)
		return
	}

	if fileMeta.BucketID != bucketMeta.ID {
		http.Error(w, "File does not belong to the specified bucket", http.StatusBadRequest)
		return
	}

	err = models.DeleteFileByID(fileId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete file: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":  "success",
		"message": "File deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}
