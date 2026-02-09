package handlers

import (
	"Cinema_System_Project/services"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type CreateSessionRequest struct {
	MovieID   string    `json:"movie_id"`
	HallName  string    `json:"hall_name"`
	StartTime time.Time `json:"start_time"`
	Price     float64   `json:"price"`
}

type CreateBookingRequest struct {
	UserID     string `json:"user_id"`
	SessionID  string `json:"session_id"`
	SeatNumber string `json:"seat_number"`
}

func CreateSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		services.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		services.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.MovieID == "" || req.HallName == "" {
		services.SendError(w, http.StatusBadRequest, "Movie ID and hall name are required")
		return
	}

	session, err := services.CreateSession(req.MovieID, req.HallName, req.StartTime, req.Price)
	if err != nil {
		services.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	services.SendSuccess(w, session)
}

func GetSessionsByMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		services.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/sessions/movie/")
	if path == "" {
		services.SendError(w, http.StatusBadRequest, "Movie ID is required")
		return
	}

	sessions, err := services.GetSessionsByMovieID(path)
	if err != nil {
		services.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	services.SendSuccess(w, sessions)
}

func CreateBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		services.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		services.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.UserID == "" || req.SessionID == "" || req.SeatNumber == "" {
		services.SendError(w, http.StatusBadRequest, "User ID, session ID, and seat number are required")
		return
	}

	booking, err := services.CreateBooking(req.UserID, req.SessionID, req.SeatNumber)
	if err != nil {
		services.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	services.SendSuccess(w, booking)
}

func GetBookingsBySession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		services.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/bookings/session/")
	if path == "" {
		services.SendError(w, http.StatusBadRequest, "Session ID is required")
		return
	}

	bookings, err := services.GetBookingsBySessionID(path)
	if err != nil {
		services.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	services.SendSuccess(w, bookings)
}

func GetUserBookings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		services.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/bookings/user/")
	if path == "" {
		services.SendError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	bookings, err := services.GetUserBookings(path)
	if err != nil {
		services.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	services.SendSuccess(w, bookings)
}

func DeleteSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		services.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/sessions/delete/")
	if path == "" {
		services.SendError(w, http.StatusBadRequest, "Session ID is required")
		return
	}

	err := services.DeleteSession(path)
	if err != nil {
		services.SendError(w, http.StatusInternalServerError, "Failed to delete session")
		return
	}

	services.SendSuccess(w, map[string]string{"message": "Session deleted successfully"})
}

func DeleteBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		services.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/bookings/delete/")
	if path == "" {
		services.SendError(w, http.StatusBadRequest, "Booking ID is required")
		return
	}

	err := services.DeleteBooking(path)
	if err != nil {
		services.SendError(w, http.StatusInternalServerError, "Failed to delete booking")
		return
	}

	services.SendSuccess(w, map[string]string{"message": "Booking deleted successfully"})
}

func GetAllBookings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		services.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	bookings, err := services.GetAllBookings()
	if err != nil {
		services.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	services.SendSuccess(w, bookings)
}
