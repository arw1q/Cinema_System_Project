package db

import (
 "Cinema_System_Project/config"
 "Cinema_System_Project/models"
 "context"
 "time"

 "go.mongodb.org/mongo-driver/bson"
 "go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllMovies() ([]models.Movie, error) {
 ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
 defer cancel()

 cursor, err := config.DB.Collection("movies").Find(ctx, bson.M{})
 if err != nil {
  return nil, err
 }
 defer cursor.Close(ctx)

 var movies []models.Movie
 if err := cursor.All(ctx, &movies); err != nil {
  return nil, err
 }
 return movies, nil
}

func InsertMovie(movie models.Movie) (primitive.ObjectID, error) {
 ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
 defer cancel()

 result, err := config.DB.Collection("movies").InsertOne(ctx, movie)
 if err != nil {
  return primitive.NilObjectID, err
 }
 return result.InsertedID.(primitive.ObjectID), nil
}

func DeleteMovie(id primitive.ObjectID) error {
 ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
 defer cancel()

 _, err := config.DB.Collection("movies").DeleteOne(ctx, bson.M{"_id": id})
 return err
}
