package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/joyboy1210/stolight/config"
	"github.com/joyboy1210/stolight/storage"
	"gorm.io/gorm"
)

type UploadHandler struct {
	Db     *gorm.DB
	Config *config.Config
}

func (u *UploadHandler) UploadHandlerAPI(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
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
	err = storage.SplitFile(u.Db, header.Filename, header.Size, file, u.Config.StorageNodes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Storage failed: %v", err), http.StatusInternalServerError)
		return
	}
	duration := time.Since(startTime)

	fileSizeMB := float64(header.Size) / (1024 * 1024)
	seconds := duration.Seconds()
	speed := fileSizeMB / seconds

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":    "success",
		"filename":  header.Filename,
		"size_mb":   fmt.Sprintf("%.2f MB", fileSizeMB),
		"time_took": fmt.Sprintf("%.2f seconds", seconds),
		"speed":     fmt.Sprintf("%.2f MB/s", speed),
	}
	json.NewEncoder(w).Encode(response)

	fmt.Printf("Upload Complete. Speed: %.2f MB/s\n", speed)
}
