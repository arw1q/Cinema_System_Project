package main

import (
	"Cinema_System_Project/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/movies", handlers.Movies)
	http.HandleFunc("/book", handlers.Book)

	http.ListenAndServe(":8080", nil)
}
