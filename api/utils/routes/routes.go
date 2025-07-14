package routes

import (
	db "api/db"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	apisecurity "api/utils/apisecurity"
	"fmt"
	"encoding/json"
	persons "api/struct"
)

// Route : /health
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Healthy")
}

// Route : /public/v1/persons
func GetPersons(w http.ResponseWriter, r *http.Request) {

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

// Route : /private/v1/guess/persons
func GetPersonsOfTheDay(w http.ResponseWriter, r *http.Request) {

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

// Route : /private/v1/guess/person/create
func CreatePersonOfTheDay(w http.ResponseWriter, r *http.Request) {

	// Check if the request is authorized with API token
	if !apisecurity.IsAuthorized(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get MongoDB client from context
	mongoClient := r.Context().Value("mongoClient").(*mongo.Client)

	// Update person of the day
	if err := db.UpdatePersonOfTheDay(mongoClient); err != nil {
		http.Error(w, "Failed to update person of the day: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Person of the day updated successfully")
}

// Route : /public/v1/guess/person/submit
func GuessPersonOfTheDay(w http.ResponseWriter, r *http.Request) {

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

// Route : /public/v1/guess/person/hint
func GetHint(w http.ResponseWriter, r *http.Request) {

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

// Route : /private/v1/guess/persons/today
func GetPersonOfTheDay(w http.ResponseWriter, r *http.Request) {

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

	if err := json.NewEncoder(w).Encode(personsGuess[len(personsGuess)-1]); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// Route : /public/v1/persons/yesterday
func GetPersonOfYesterday(w http.ResponseWriter, r *http.Request) {

	// Get MongoDB client from context
	mongoClient := r.Context().Value("mongoClient").(*mongo.Client)

	// Get person of yesterday
	personOfYesterday, err := db.GetPersonOfYesterday(mongoClient, "dodle")
	if err != nil {
		http.Error(w, "Failed to get person of yesterday: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(personOfYesterday); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// Route : /public/v1/guess/id
func GetGuessID(w http.ResponseWriter, r *http.Request) {

	// Get MongoDB client from context
	mongoClient := r.Context().Value("mongoClient").(*mongo.Client)

	id, error := db.GetGuessID(mongoClient, "dodle")

	if id == "" || error != nil {
		http.Error(w, "No guess found for today", http.StatusNotFound)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(map[string]string{"id": id}); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}