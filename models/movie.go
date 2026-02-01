package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Movie struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title    string             `bson:"title"         json:"title"`
	Genre    string             `bson:"genre"         json:"genre"`
	Year     int                `bson:"year"          json:"year"`
	Director string             `bson:"director"      json:"director"`
}
