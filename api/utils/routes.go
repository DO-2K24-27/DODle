package routes

import (
	"api/db"
	"api/utils/api_security"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/health" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Healthy")
}

func GetPersons(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/public/v1/persons" {
		http.NotFound(w, r)
		return
	}

	// Get MongoDB client from context
	mongoClient := r.Context().Value("mongoClient").(*mongo.Client)

	// Get persons from the database
	personsList, err := db.GetPersons(mongoClient, "dodle")
	if err != nil {
		http.Error(w, "Failed to get persons: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Encode persons as JSON and send response
	if err := json.NewEncoder(w).Encode(personsList); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}

func GetPersonsOfTheDay(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/private/v1/guess/persons" {
		http.NotFound(w, r)
		return
	}

	// Check if the request is authorized with API token
	if !apisecurity.IsAuthorized(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get MongoDB client from context
	mongoClient := r.Context().Value("mongoClient").(*mongo.Client)

	// Update person of the day
	personsGuess, err := db.GetPersonsOfTheDay(mongoClient, "dodle")
	if err != nil {
		http.Error(w, "Failed to update person of the day: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(personsGuess); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreatePersonOfTheDay(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/private/v1/guess/person/create" {
		http.NotFound(w, r)
		return
	}

	// Check if the request is authorized with API token
	if !apisecurity.IsAuthorized(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get MongoDB client from context
	mongoClient := r.Context().Value("mongoClient").(*mongo.Client)

	// Update person of the day
	if err := UpdatePersonOfTheDay(mongoClient); err != nil {
		http.Error(w, "Failed to update person of the day: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Person of the day updated successfully")
}

func GuessPersonOfTheDay(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/public/v1/guess/person/submit" {
		http.NotFound(w, r)
		return
	}

	// Get MongoDB client from context
	mongoClient := r.Context().Value("mongoClient").(*mongo.Client)

	// Decode the guess from the request body
	var guess persons.Person
	if err := json.NewDecoder(r.Body).Decode(&guess); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Received guess:", guess)

	// Try to guess the person of the day
	correct, returnedPerson, err := db.TryGuess(mongoClient, "dodle", guess)
	if err != nil {
		http.Error(w, "Failed to process guess: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"correct": correct,
		"person":  returnedPerson,
	}); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetHint(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/public/v1/guess/person/hint" {
		http.NotFound(w, r)
		return
	}

	// Get MongoDB client from context
	mongoClient := r.Context().Value("mongoClient").(*mongo.Client)

	// Get person of the day
	personOfTheDay, err := db.GetPersonOfTheDay(mongoClient, "dodle")
	if err != nil {
		http.Error(w, "Failed to get person of the day: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(personOfTheDay.Hint); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}