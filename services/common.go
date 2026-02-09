package services

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Token   string      `json:"token,omitempty"`
}

func SendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func SendError(w http.ResponseWriter, statusCode int, message string) {
	SendJSON(w, statusCode, Response{
		Success: false,
		Message: message,
	})
}

func SendSuccess(w http.ResponseWriter, data interface{}) {
	SendJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func SendSuccessWithToken(w http.ResponseWriter, data interface{}, token string) {
	SendJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
		Token:   token,
	})
}
