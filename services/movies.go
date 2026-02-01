package services

import (
	"Cinema_System_Project/db"
	"Cinema_System_Project/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetMovies fetches all movies using a goroutine + channel (concurrency requirement).
func GetMovies() ([]models.Movie, error) {
	type result struct {
		movies []models.Movie
		err    error
	}

	ch := make(chan result, 1)

	go func() {
		movies, err := db.GetAllMovies()
		ch <- result{movies, err}
	}()

	res := <-ch
	return res.movies, res.err
}

func AddMovie(movie models.Movie) (primitive.ObjectID, error) {
	return db.InsertMovie(movie)
}

func RemoveMovie(id primitive.ObjectID) error {
	return db.DeleteMovie(id)
}
