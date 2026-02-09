package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username     string             `json:"username" bson:"username"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"-" bson:"password_hash"`
	Role         string             `json:"role" bson:"role"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}

type Movie struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Duration    int                `json:"duration" bson:"duration"`
	PosterURL   string             `json:"poster_url" bson:"poster_url"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type Session struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	MovieID   primitive.ObjectID `json:"movie_id" bson:"movie_id"`
	HallName  string             `json:"hall_name" bson:"hall_name"`
	StartTime time.Time          `json:"start_time" bson:"start_time"`
	Price     float64            `json:"price" bson:"price"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type Booking struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"`
	SessionID  primitive.ObjectID `json:"session_id" bson:"session_id"`
	SeatNumber string             `json:"seat_number" bson:"seat_number"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
