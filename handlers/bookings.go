package handlers

import "net/http"

func Book(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("booking created"))
}
