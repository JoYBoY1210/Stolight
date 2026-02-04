package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joyboy1210/stolight/config"
	"github.com/joyboy1210/stolight/storage"
)

func UploadHandlerAPI(w http.ResponseWriter, r *http.Request) {

	bucketName := r.URL.Query().Get("bucket")
	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()
	nodes := config.Cfg.StorageNodes
	err = storage.SplitFile(header.Filename, header.Size, file, nodes, bucketName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Storage failed: %v", err), http.StatusInternalServerError)
		return
	}

	fileSizeMB := float64(header.Size) / (1024 * 1024)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":   "success",
		"filename": header.Filename,
		"size_mb":  fmt.Sprintf("%.2f MB", fileSizeMB),
		"bucket":   bucketName,
	}
	json.NewEncoder(w).Encode(response)

}
