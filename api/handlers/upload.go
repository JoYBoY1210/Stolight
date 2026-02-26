package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/joyboy1210/stolight/queue"
	"github.com/joyboy1210/stolight/storage"
)

func UploadHandlerAPI(w http.ResponseWriter, r *http.Request) {

	bucketName := r.PathValue("bucket")
	if bucketName == "" {
		http.Error(w, "Bucket name is required in the URL", http.StatusBadRequest)
		return
	}

	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "Failed to read multipart data", http.StatusBadRequest)
		return
	}

	var fileID string
	var size int64
	var fileName string
	var fileFound bool

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Failed to read multipart data", http.StatusBadRequest)
			return
		}

		if part.FormName() == "file" {
			fileFound = true
			fileName = part.FileName()

			fileID, size, err = storage.StageFile(part, fileName, 0, bucketName)
			part.Close()

			if err != nil {
				http.Error(w, fmt.Sprintf("Storage failed: %v", err), http.StatusInternalServerError)
				return
			}
			q := queue.GetQueue()
			err = q.AddJob(fileID)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to add job to queue: %v", err), http.StatusInternalServerError)
				return
			}
			break
		}
		part.Close()
	}

	if !fileFound {
		http.Error(w, "Expected form field 'file'", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":   "processing",
		"filename": fileName,
		"file_id":  fileID,
		"size":     size,
		"bucket":   bucketName,
	}
	json.NewEncoder(w).Encode(response)

}
