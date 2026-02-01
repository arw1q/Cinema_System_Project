package handlers

import (
	"Cinema_System_Project/models"
	"Cinema_System_Project/services"
	"encoding/json"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Movies(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getMovies(w, r)
	case http.MethodPost:
		addMovie(w, r)
	case http.MethodDelete:
		removeMovie(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /movies — any user can browse
func getMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := services.GetMovies()
	if err != nil {
		http.Error(w, "failed to fetch movies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// POST /movies?role=admin — admin only
func addMovie(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("role") != "admin" {
		http.Error(w, "forbidden: admin only", http.StatusForbidden)
		return
	}

	var movie models.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := services.AddMovie(movie)
	if err != nil {
		http.Error(w, "failed to add movie", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id.Hex()})
}

// DELETE /movies/{id}?role=admin — admin only
func removeMovie(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("role") != "admin" {
		http.Error(w, "forbidden: admin only", http.StatusForbidden)
		return
	}

	// extract the id from the end of the path: /movies/{id}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	idStr := parts[len(parts)-1]

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "invalid movie id", http.StatusBadRequest)
		return
	}

	if err := services.RemoveMovie(id); err != nil {
		http.Error(w, "failed to delete movie", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "movie deleted"})
}
