package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/joyboy1210/stolight/models"
	"github.com/joyboy1210/stolight/storage"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileId := strings.TrimPrefix(r.URL.Path, "/download/")
	if fileId == "" {
		http.Error(w, "file ID is required", http.StatusBadRequest)
		return
	}
	fileMeta, err := models.GetFileByID(fileId)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileMeta.Name))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileMeta.Size))

	err = storage.MergeFile(fileId, w)
	if err != nil {
		http.Error(w, "failed to download file", http.StatusInternalServerError)
		return
	}
}
