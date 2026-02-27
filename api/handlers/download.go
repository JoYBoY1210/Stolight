package handlers

import (
	"fmt"
	"net/http"

	"github.com/joyboy1210/stolight/config"
	"github.com/joyboy1210/stolight/models"
	"github.com/joyboy1210/stolight/storage"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bucketName := r.PathValue("bucket")
	fileId := r.PathValue("fileId") //i need to change this handler to use fileId instead of fileName, but for now i will keep it as fileName to make it work with current implementation

	if bucketName == "" {
		http.Error(w, "bucket name is required", http.StatusBadRequest)
		return
	}

	if fileId == "" {
		http.Error(w, "file Id is required", http.StatusBadRequest)
		return
	}

	fileMeta, err := models.GetFileByID(fileId)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	fileName := fileMeta.Name
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	// w.Header().Set("Content-Length", fmt.Sprintf("%d", fileMeta.Size))

	nodes := config.Cfg.StorageNodes

	err = storage.DecodeFile(w, fileId, nodes, fileMeta.Size)
	if err != nil {
		fmt.Println("failed to download file:", err)
		return
	}
}
