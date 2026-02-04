package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joyboy1210/stolight/models"
)

func CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid name err:- %s", err.Error()), http.StatusBadRequest)
	}
	if req.Name == "" {
		http.Error(w, "Bucket name cannot be empty", http.StatusBadRequest)
		return
	}
	bucket, err := models.CreateBucket(req.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create bucket: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"status": "success",
		"bucket": bucket,
	}
	json.NewEncoder(w).Encode(response)
}


