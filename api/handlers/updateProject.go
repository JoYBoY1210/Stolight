package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/joyboy1210/stolight/models"
)

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("sto-Key")
	if key == "" {
		http.Error(w, "Only admin can access this", http.StatusUnauthorized)
		return
	}
	user, err := models.GetUserByKey(key)
	if err != nil || user.Role != "admin" {
		http.Error(w, "Only admin can access this", http.StatusUnauthorized)
		return
	}
	var req struct {
		Username       string `json:"username"`
		AllowedBuckets string `json:"allowed_buckets"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Username == "" {
		http.Error(w, "Project name is required", http.StatusBadRequest)
		return
	}
	project, err := models.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	if project.Role == "admin" {
		http.Error(w, "Cannot update admin permissions", http.StatusBadRequest)
		return
	}

	project.AllowedBuckets = req.AllowedBuckets
	if err := models.UpdateUser(project); err != nil {
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"msg":      "Permissions updated",
		"username": project.Username,
		"buckets":  project.AllowedBuckets,
	})
}
