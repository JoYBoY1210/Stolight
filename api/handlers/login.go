package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/joyboy1210/stolight/models"
	"github.com/joyboy1210/stolight/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := models.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if user.PasswordHash == "" {
		http.Error(w, "login not allowed for service accounts", http.StatusForbidden)
		return
	}
	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"key":  user.Key,
		"role": user.Role,
	})
}
