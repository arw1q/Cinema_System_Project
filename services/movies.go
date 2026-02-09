package services

import (
	"Cinema_System_Project/db"
	"Cinema_System_Project/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateMovie(title, description, posterURL string, duration int) (*models.Movie, error) {
	collection := db.GetCollection("movies")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	movie := &models.Movie{
		ID:          primitive.NewObjectID(),
		Title:       title,
		Description: description,
		Duration:    duration,
		PosterURL:   posterURL,
		CreatedAt:   time.Now(),
	}

	_, err := collection.InsertOne(ctx, movie)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func GetAllMovies() ([]models.Movie, error) {
	collection := db.GetCollection("movies")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movies []models.Movie
	if err = cursor.All(ctx, &movies); err != nil {
		return nil, err
	}

	return movies, nil
}

func GetMovieByID(id string) (*models.Movie, error) {
	collection := db.GetCollection("movies")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var movie models.Movie
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&movie)
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func DeleteMovie(id string) error {
	collection := db.GetCollection("movies")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
