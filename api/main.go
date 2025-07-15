package main

import (
	db "api/db"
	ctxUtil "api/utils/context"
	routes "api/utils/routes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

// Middleware to inject MongoDB client into request context
func withMongoClient(next http.Handler, client *mongo.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxUtil.MongoClientKey, client)
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

		// Override MongoDB URI environment variable
		if err := os.Setenv("MONGODB_URI", mongoURI); err != nil {
			log.Printf("Warning: Failed to set MONGODB_URI environment variable: %v", err)
		}
	}

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
	if err := db.InitDB(mongoClient); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create router
	mux := http.NewServeMux()

	// Wrap handlers with MongoDB client middleware
	personsHandler := http.HandlerFunc(routes.GetPersons)
	guessesHandler := http.HandlerFunc(routes.GetPersonsOfTheDay)
	guessHandler := http.HandlerFunc(routes.GetPersonOfTheDay)
	createPODHandler := http.HandlerFunc(routes.CreatePersonOfTheDay)
	guessPersonHandler := http.HandlerFunc(routes.GuessPersonOfTheDay)
	getHintHandler := http.HandlerFunc(routes.GetHint)
	getYesterdayHandler := http.HandlerFunc(routes.GetPersonOfYesterday)
	GetGuessIDHandler := http.HandlerFunc(routes.GetGuessID)

	// Register handlers
	mux.HandleFunc("/health", routes.HealthHandler)
	mux.Handle("/public/v1/persons", withMongoClient(personsHandler, mongoClient))
	mux.Handle("/public/v1/guess/person/submit", withMongoClient(guessPersonHandler, mongoClient))
	mux.Handle("/public/v1/guess/person/hint", withMongoClient(getHintHandler, mongoClient))
	mux.Handle("/public/v1/guess/person/yesterday", withMongoClient(getYesterdayHandler, mongoClient))
	mux.Handle("/public/v1/guess/id", withMongoClient(GetGuessIDHandler, mongoClient))

	mux.Handle("/private/v1/guess/persons", withMongoClient(guessesHandler, mongoClient))
	mux.Handle("/private/v1/guess/person/today", withMongoClient(guessHandler, mongoClient))
	mux.Handle("/private/v1/guess/person/create", withMongoClient(createPODHandler, mongoClient))

	fmt.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
