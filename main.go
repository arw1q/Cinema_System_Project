package main

import (
	"Cinema_System_Project/config"
	"Cinema_System_Project/handlers"
	"log"
	"net/http"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	config.ConnectDB()

	http.HandleFunc("/login", corsMiddleware(handlers.Login))
	http.HandleFunc("/movies/", corsMiddleware(handlers.Movies))
	http.HandleFunc("/bookings", corsMiddleware(handlers.Book))

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
