package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/joyboy1210/stolight/models"
)

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("sto-Key")
	if key == "" {
		http.Error(w, "Missing API key", http.StatusUnauthorized)
		return
	}
	user, err := models.GetUserByKey(key)
	if err != nil {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	}
	if user.Role != "admin" {
		http.Error(w, "Forbidden: Only admin users can create projects", http.StatusForbidden)
		return
	}
	var req struct {
		Name    string `json:"name"`
		Buckets string `json:"buckets"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "Project name is required", http.StatusBadRequest)
		return
	}
	apiKey := uuid.New().String()
	ID := uuid.New().String()
	err = models.CreateUser(ID, req.Name, apiKey, "project", req.Buckets)
	if err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"id":     ID,
		"key":    apiKey,
	})
}
