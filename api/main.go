package main

import (
	db "api/db"
	persons "api/struct"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"math/rand"

	"go.mongodb.org/mongo-driver/mongo"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Hello world!")
}

func getPersons(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/persons" {
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
	if r.URL.Path != "/guess/persons" {
		http.NotFound(w, r)
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
	if r.URL.Path != "/guess/person/create" {
		http.NotFound(w, r)
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

func initDB(mongoClient *mongo.Client) error {
	db.CreateDatabase(mongoClient, "dodle")
	fmt.Println("Creating collections...")
	db.CreateCollection(mongoClient, "dodle", "Persons")
	db.CreateCollection(mongoClient, "dodle", "GuessesOfTheMonth")

	// Load persons from file
	persons, err := OpenPersonsFile()
	if err != nil {
		return fmt.Errorf("failed to open persons file: %v", err)
	}

	// Check if collection is empty before populating
	existingPersons, err := db.GetPersons(mongoClient, "dodle")
	if err != nil {
		return fmt.Errorf("failed to check existing persons: %v", err)
	}

	if len(existingPersons.Persons) == 0 {
		fmt.Println("Populating Persons collection...")
		result := db.PopulatePersonsCollection(mongoClient, "dodle", persons)
		if result != "" {
			return fmt.Errorf("failed to populate persons collection: %s", result)
		}
	} else {
		fmt.Println("Persons collection already contains data, skipping population")
	}

	return nil
}

func OpenPersonsFile() (persons.Persons, error) {
	jsonFile, err := os.Open("./data/persons.json")
	if err != nil {
		return persons.Persons{}, fmt.Errorf("error opening file: %v", err)
	}
	defer jsonFile.Close()

	// Read the content of the file
	content, err := io.ReadAll(jsonFile)
	if err != nil {
		return persons.Persons{}, fmt.Errorf("error reading file: %v", err)
	}

	// Unmarshal the JSON data into the struct
	var personsList []persons.Person
	err = json.Unmarshal(content, &personsList)
	if err != nil {
		return persons.Persons{}, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return persons.Persons{Persons: personsList}, nil
}

func UpdatePersonOfTheDay(mongoClient *mongo.Client) error {
	previousPersons, err := db.GetPersonsOfTheDay(mongoClient, "dodle")
	if err != nil {
		return fmt.Errorf("failed to get previous persons of the day: %v", err)
	}

	personsAvailable, err := db.GetPersons(mongoClient, "dodle")
	if err != nil {
		return fmt.Errorf("failed to get persons: %v", err)
	}

	if len(personsAvailable.Persons) == 0 {
		return fmt.Errorf("no persons available to update person of the day")
	}

	isSelectable := false
	candidate := personsAvailable.Persons[rand.Intn(len(personsAvailable.Persons))]
	fmt.Println("New candidate for person of the day:", candidate.Firstname, candidate.Lastname)
	for i := 0; i < len(personsAvailable.Persons); i++ {
		candidate = personsAvailable.Persons[rand.Intn(len(personsAvailable.Persons))]
		isSelectable = true
		for _, person := range previousPersons {
			if person.Firstname == candidate.Firstname && person.Lastname == candidate.Lastname {
				isSelectable = false
			}
		}

		if isSelectable {
			break
		}
	}

	dateOfToday := time.Now().Format("2006-01-02")

	// If we already selected a candidate today we delete the previous one
	if GetPersonOfTheDay, err := db.GetPersonOfTheDay(mongoClient, "dodle"); err == nil {
		if GetPersonOfTheDay.Firstname != "" {
			if err := db.DeletePersonOfTheDay(mongoClient, "dodle", dateOfToday); err != nil {
				return fmt.Errorf("failed to delete previous person of the day: %v", err)
			}
			fmt.Println("Previous person of the day deleted successfully")
		}
	}

	if err := db.CreatePersonOfTheDay(mongoClient, "dodle", candidate); err != nil {
		return fmt.Errorf("failed to create person of the day: %v", err)
	}

	dateToDelete := time.Now().AddDate(0, 0, -10).Format("2006-01-02")

	if err := db.DeletePersonOfTheDay(mongoClient, "dodle", dateToDelete); err != nil {
		return fmt.Errorf("failed to delete previous person of the day: %v", err)
	}
	fmt.Println("Previous person of the day deleted successfully")
	// You can add your logic here to update the person of the day in the database.
	return nil
}

func GuessPersonOfTheDay(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/guess/person/submit" {
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
	if r.URL.Path != "/guess/person/hint" {
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

// Middleware to inject MongoDB client into request context
func withMongoClient(next http.Handler, client *mongo.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "mongoClient", client)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	fmt.Println("Starting the server...")

	// Set MongoDB URI from environment variable or use default
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		// Default URI with authentication
		mongoURI = "mongodb://admin:admin@localhost:27017"
	}

	// Override MongoDB URI environment variable
	os.Setenv("MONGODB_URI", mongoURI)

	// Connect to MongoDB
	mongoClient, err := db.ConnectToMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ensure MongoDB client is closed when the program exits
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	fmt.Println("Connected to MongoDB successfully!")

	// Initialize database
	if err := initDB(mongoClient); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create router
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/hello", helloHandler)

	// Wrap handlers with MongoDB client middleware
	personsHandler := http.HandlerFunc(getPersons)
	guessHandler := http.HandlerFunc(GetPersonsOfTheDay)
	createPODHandler := http.HandlerFunc(CreatePersonOfTheDay)
	guessPersonHandler := http.HandlerFunc(GuessPersonOfTheDay)
	getHintHandler := http.HandlerFunc(GetHint)
	mux.Handle("/persons", withMongoClient(personsHandler, mongoClient))
	mux.Handle("/guess/persons", withMongoClient(guessHandler, mongoClient))
	mux.Handle("/guess/person/create", withMongoClient(createPODHandler, mongoClient))
	mux.Handle("/guess/person/submit", withMongoClient(guessPersonHandler, mongoClient))
	mux.Handle("/guess/person/hint", withMongoClient(getHintHandler, mongoClient))

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
