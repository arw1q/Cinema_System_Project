package main

import (
	"Cinema_System_Project/db"
	"Cinema_System_Project/handlers"
	"Cinema_System_Project/middleware"
	"log"
	"net/http"
	"os"
	"strings"
)

func loadEnv() {
	data, err := os.ReadFile(".env")
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			os.Setenv(parts[0], parts[1])
		}
	}
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func handleMoviesWithID(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/movies/") && r.URL.Path != "/api/movies/" {
		handlers.GetMovie(w, r)
	} else {
		handlers.GetMovies(w, r)
	}
}

func handleSessionsWithMovie(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/sessions/movie/") {
		handlers.GetSessionsByMovie(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func handleBookingsWithID(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/bookings/session/") {
		handlers.GetBookingsBySession(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/api/bookings/user/") {
		handlers.GetUserBookings(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func main() {
	loadEnv()
	db.Connect()

	http.HandleFunc("/api/auth/register", enableCORS(handlers.Register))
	http.HandleFunc("/api/auth/login", enableCORS(handlers.Login))
	http.HandleFunc("/api/auth/admin", enableCORS(handlers.AdminLogin))

	http.HandleFunc("/api/movies/create", enableCORS(middleware.AdminMiddleware(handlers.CreateMovie)))
	http.HandleFunc("/api/movies", enableCORS(handleMoviesWithID))
	http.HandleFunc("/api/movies/", enableCORS(handleMoviesWithID))
	http.HandleFunc("/api/movies/delete/", enableCORS(middleware.AdminMiddleware(handlers.DeleteMovie)))

	http.HandleFunc("/api/sessions/create", enableCORS(middleware.AdminMiddleware(handlers.CreateSession)))
	http.HandleFunc("/api/sessions/", enableCORS(handleSessionsWithMovie))
	http.HandleFunc("/api/sessions/delete/", enableCORS(middleware.AdminMiddleware(handlers.DeleteSession)))

	http.HandleFunc("/api/bookings/create", enableCORS(middleware.AuthMiddleware(handlers.CreateBooking)))
	http.HandleFunc("/api/bookings/", enableCORS(handleBookingsWithID))
	http.HandleFunc("/api/bookings/delete/", enableCORS(middleware.AdminMiddleware(handlers.DeleteBooking)))
	http.HandleFunc("/api/bookings/all", enableCORS(middleware.AdminMiddleware(handlers.GetAllBookings)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server starting on port %s...", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
