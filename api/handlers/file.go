package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/joyboy1210/stolight/models"
)

func ListFilesInBucketHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	prefix := "/api/files/all/"
	if strings.HasPrefix(r.URL.Path, prefix) == false {
		http.Error(w, "Invalid URL to list files in bucket", http.StatusBadRequest)
		return
	}
	bucketName := strings.TrimPrefix(r.URL.Path, prefix)
	if bucketName == "" {
		http.Error(w, "Bucket name is required to list all the files in it", http.StatusBadRequest)
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
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	prefix := "/api/files/delete/"
	if strings.HasPrefix(r.URL.Path, prefix) == false {
		http.Error(w, "Invalid URL to delete file", http.StatusBadRequest)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, prefix)
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		http.Error(w, "Invalid download URL", http.StatusBadRequest)
		return
	}
	bucketName := parts[0]
	fileName := parts[1]

	if bucketName == "" {
		http.Error(w, "bucket name is required", http.StatusBadRequest)
		return
	}

	if fileName == "" {
		http.Error(w, "file name is required", http.StatusBadRequest)
		return
	}

	fileName, err := url.PathUnescape(fileName)
	if err != nil {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	bucketName, err = url.PathUnescape(bucketName)
	if err != nil {
		http.Error(w, "Invalid bucket name", http.StatusBadRequest)
		return
	}

	if strings.Contains(fileName, "..") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	bucket, err := models.GetBucketByName(bucketName)
	if err != nil {
		http.Error(w, "bucket not found", http.StatusNotFound)
		return
	}

	fileMeta, err := models.GetFileByFileNameAndBucketId(fileName, bucket.ID)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	fileID := fileMeta.ID

	err = models.DeleteFileByID(fileID)
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
