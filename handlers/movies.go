package handlers

import (
	"Cinema_System_Project/services"
	"encoding/json"
	"net/http"
	"strings"
)

type CreateMovieRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	PosterURL   string `json:"poster_url"`
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		services.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req CreateMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		services.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" || req.Duration <= 0 {
		services.SendError(w, http.StatusBadRequest, "Title and duration are required")
		return
	}

	movie, err := services.CreateMovie(req.Title, req.Description, req.PosterURL, req.Duration)
	if err != nil {
		services.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	services.SendSuccess(w, movie)
}

func GetMovies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		services.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	movies, err := services.GetAllMovies()
	if err != nil {
		services.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	services.SendSuccess(w, movies)
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		services.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/movies/")
	if path == "" {
		services.SendError(w, http.StatusBadRequest, "Movie ID is required")
		return
	}

	movie, err := services.GetMovieByID(path)
	if err != nil {
		services.SendError(w, http.StatusNotFound, "Movie not found")
		return
	}

	services.SendSuccess(w, movie)
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		services.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/movies/delete/")
	if path == "" {
		services.SendError(w, http.StatusBadRequest, "Movie ID is required")
		return
	}

	err := services.DeleteMovie(path)
	if err != nil {
		services.SendError(w, http.StatusInternalServerError, "Failed to delete movie")
		return
	}

	services.SendSuccess(w, map[string]string{"message": "Movie deleted successfully"})
}
