# Baha Cinema

Baha Cinema is an Advanced Programming Final Project that represents a cinema management system. The project demonstrates backend service development, database integration, and basic frontend interaction.

The system provides a REST API for user authentication, movie management, session scheduling, and ticket booking. The backend is written in Go and follows a structured, layered architecture.

## Technologies
Go, REST API, Database, HTML, CSS, JavaScript

## Running the Project
Clone the repository, install dependencies using `go mod download`, configure the database connection, and start the server with `go run main.go`.

## API Endpoints
POST   /api/auth/register  
POST   /api/auth/login  
POST   /api/auth/admin  

POST   /api/movies/create (Auth required)  
GET    /api/movies  
DELETE /api/movies/delete/{id} (Auth required)  

POST   /api/sessions/create (Auth required)  
GET    /api/sessions/movie/{id}  
DELETE /api/sessions/delete/{id} (Auth required)  

POST   /api/bookings/create (Auth required)  
GET    /api/bookings/all (Auth required)

## Project Type
Advanced Programming Final Project

## Authors
Arlan  
Ersultan  
Bakyt
