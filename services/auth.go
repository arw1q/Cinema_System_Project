package services

import (
	"Cinema_System_Project/db"
	"Cinema_System_Project/models"
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func GenerateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.ID.Hex(),
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func RegisterUser(username, email, password string) (*models.User, string, error) {
	collection := db.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"$or": []bson.M{
		{"username": username},
		{"email": email},
	}}).Decode(&existingUser)
	if err == nil {
		return nil, "", errors.New("username or email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		ID:           primitive.NewObjectID(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "user",
		CreatedAt:    time.Now(),
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return nil, "", err
	}

	token, err := GenerateJWT(*user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func LoginUser(username, password string) (*models.User, string, error) {
	collection := db.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	token, err := GenerateJWT(user)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func LoginAdmin(username, password string) (string, error) {
	if username != "admin" || password != "12345" {
		return "", errors.New("invalid admin credentials")
	}

	adminUser := models.User{
		ID:       primitive.NewObjectID(),
		Username: "admin",
		Email:    "admin@cinema.com",
		Role:     "admin",
	}

	token, err := GenerateJWT(adminUser)
	if err != nil {
		return "", err
	}

	return token, nil
}
