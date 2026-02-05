package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/joyboy1210/stolight/models"
	"github.com/joyboy1210/stolight/storage"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	prefix := "/api/download/"
	if strings.HasPrefix(r.URL.Path, prefix) == false {
		http.Error(w, "Invalid download URL", http.StatusBadRequest)
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
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileMeta.Name))
	w.Header().Set("Content-Type", "application/octet-stream")
	// w.Header().Set("Content-Length", fmt.Sprintf("%d", fileMeta.Size))

	err = storage.MergeFile(fileMeta.ID, w)
	if err != nil {
		fmt.Println("failed to download file:", err)
		return
	}
}
