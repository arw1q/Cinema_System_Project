package handlers

import (
	"Cinema_System_Project/models"
	"Cinema_System_Project/services"
	"encoding/json"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		services.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		services.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		services.SendError(w, http.StatusBadRequest, "Username, email and password are required")
		return
	}

	user, token, err := services.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		services.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	services.SendSuccessWithToken(w, user, token)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		services.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		services.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		services.SendError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	user, token, err := services.LoginUser(req.Username, req.Password)
	if err != nil {
		services.SendError(w, http.StatusUnauthorized, err.Error())
		return
	}

	services.SendSuccessWithToken(w, user, token)
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		services.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		services.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, err := services.LoginAdmin(req.Username, req.Password)
	if err != nil {
		services.SendError(w, http.StatusUnauthorized, "Invalid admin credentials")
		return
	}

	services.SendSuccessWithToken(w, map[string]string{
		"role":     "admin",
		"username": "admin",
		"email":    "admin@cinema.com",
	}, token)
}
