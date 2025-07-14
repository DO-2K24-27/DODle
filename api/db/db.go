package db

import (
	persons "api/struct"
	data "api/utils/data"
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"math/rand"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB() (*mongo.Client, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return client, nil
}

func CreateDatabase(client *mongo.Client, dbName string) string {
	// Create a database by accessing it
	db := client.Database(dbName)
	if db == nil {
		return "Failed to create database"
	}
	fmt.Printf("Database %s created successfully\n", dbName)
	return ""
}

func CreateCollection(client *mongo.Client, dbName, collectionName string) string {
	// Create a collection in the specified database
	collection := client.Database(dbName).Collection(collectionName)
	if collection == nil {
		return "Failed to create collection"
	}
	fmt.Printf("Collection %s created in database %s", collectionName, dbName)
	return ""
}

func PopulatePersonsCollection(client *mongo.Client, dbName string, personsToInsert persons.Persons) string {
	// Example function to populate the Persons collection
	// This is a placeholder; actual implementation would depend on your data model
	collection := client.Database(dbName).Collection("Persons")
	if collection == nil {
		return "Failed to access Persons collection"
	}

	// Insert persons into the collection
	personsInterface := []interface{}{}
	for _, person := range personsToInsert.Persons {
		personsInterface = append(personsInterface, person)
	}
	if len(personsInterface) == 0 {
		return "No persons to insert"
	}
	_, err := collection.InsertMany(context.TODO(), personsInterface)
	if err != nil {
		return "Failed to insert persons: " + err.Error()
	}
	fmt.Printf("Inserted %d persons into the Persons collection", len(personsInterface))
	return ""
}

func GetPersons(client *mongo.Client, dbName string) (persons.Persons, error) {
	// Retrieve all persons from the Persons collection
	collection := client.Database(dbName).Collection("Persons")
	if collection == nil {
		return persons.Persons{}, nil
	}

	cursor, err := collection.Find(context.TODO(), map[string]interface{}{})
	if err != nil {
		return persons.Persons{}, err
	}
	defer cursor.Close(context.TODO())

	var personsList persons.Persons
	for cursor.Next(context.TODO()) {
		var person persons.Person
		if err := cursor.Decode(&person); err != nil {
			return persons.Persons{}, err
		}
		personsList.Persons = append(personsList.Persons, person)
	}

	if err := cursor.Err(); err != nil {
		return persons.Persons{}, err
	}

	return personsList, nil
}

func GetPersonOfTheDay(client *mongo.Client, dbName string) (persons.Person, error) {
	// Retrieve the person of the day from the GuessesOfTheMonth collection
	collection := client.Database(dbName).Collection("GuessesOfTheMonth")
	if collection == nil {
		return persons.Person{}, fmt.Errorf("GuessesOfTheMonth collection not found")
	}

	var currentDate string = time.Now().Format("2006-01-02") // Format the date as YYYY-MM-DD

	// Create a document structure to match what's stored in MongoDB
	var doc struct {
		Date   string         `bson:"date"`
		Person persons.Person `bson:"person"`
	}

	err := collection.FindOne(context.TODO(), map[string]interface{}{"date": currentDate}).Decode(&doc)
	if err != nil {
		return persons.Person{}, fmt.Errorf("failed to find person of the day: %v", err)
	}

	// Log the found person for debugging
	fmt.Printf("Found Person of the Day (%s): %s %s\n", doc.Date, doc.Person.Firstname, doc.Person.Lastname)

	return doc.Person, nil
}

func GetPersonsOfTheDay(client *mongo.Client, dbName string) ([]persons.Person, error) {
	// Retrieve all persons of the day from the GuessesOfTheMonth collection
	collection := client.Database(dbName).Collection("GuessesOfTheMonth")
	if collection == nil {
		return nil, fmt.Errorf("GuessesOfTheMonth collection not found")
	}

	cursor, err := collection.Find(context.TODO(), map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to find persons of the day: %v", err)
	}
	defer cursor.Close(context.TODO())

	var personsOfTheDay []persons.Person
	for cursor.Next(context.TODO()) {
		// Create a document to hold the raw document structure
		var doc struct {
			Date   string         `bson:"date"`
			Person persons.Person `bson:"person"`
		}

		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode document: %v", err)
		}

		// Only add the person if it has data
		if doc.Person.Firstname != "" && doc.Person.Lastname != "" {
			fmt.Printf("Person of the Day (%s): %s %s\n", doc.Date, doc.Person.Firstname, doc.Person.Lastname)
			personsOfTheDay = append(personsOfTheDay, doc.Person)
		} else {
			fmt.Printf("Warning: Found document for date %s but person data is incomplete\n", doc.Date)
		}
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return personsOfTheDay, nil
}

func TryGuess(client *mongo.Client, dbName string, guess persons.Person) (bool, persons.Person, error) {
	// Check if the guess matches the person of the day
	personOfTheDay, err := GetPersonOfTheDay(client, dbName)
	if err != nil {
		return false, persons.Person{}, fmt.Errorf("failed to get person of the day: %v", err)
	}

	fmt.Printf("Guess: %s %s, Person of the Day: %s %s\n", guess.Firstname, guess.Lastname, personOfTheDay.Firstname, personOfTheDay.Lastname)

	// Compare the guess with the person of the day
	if guess.Firstname == personOfTheDay.Firstname && guess.Lastname == personOfTheDay.Lastname {
		return true, personOfTheDay, nil // Correct guess
	}

	var returnedPerson persons.Person
	if personOfTheDay.Firstname == guess.Firstname {
		returnedPerson.Firstname = personOfTheDay.Firstname
	}
	if personOfTheDay.Lastname == guess.Lastname {
		returnedPerson.Lastname = personOfTheDay.Lastname
	}
	if personOfTheDay.Gender == guess.Gender {
		returnedPerson.Gender = personOfTheDay.Gender
	}
	if personOfTheDay.Workplace == guess.Workplace {
		returnedPerson.Workplace = personOfTheDay.Workplace
	}
	if personOfTheDay.Type == guess.Type {
		returnedPerson.Type = personOfTheDay.Type
	}
	if personOfTheDay.Image == guess.Image {
		returnedPerson.Image = personOfTheDay.Image
	}
	if personOfTheDay.Hint == guess.Hint {
		returnedPerson.Hint = personOfTheDay.Hint
	}

	return false, returnedPerson, nil // Incorrect guess
}

func CreatePersonOfTheDay(client *mongo.Client, dbName string, person persons.Person) error {
	// Create or update the person of the day in the Persons collection
	collection := client.Database(dbName).Collection("GuessesOfTheMonth")
	if collection == nil {
		return fmt.Errorf("persons collection not found")
	}

	currentDate := time.Now().Format("2006-01-02") // Format the date as YYYY-MM-DD

	personOfTheDay := map[string]interface{}{
		"date":   currentDate,
		"person": person,
	}

	// Upsert the person of the day
	_, err := collection.InsertMany(
		context.TODO(),
		[]interface{}{personOfTheDay},
	)
	if err != nil {
		return fmt.Errorf("failed to create or update person of the day: %v", err)
	}

	return nil
}

func DeletePersonOfTheDay(client *mongo.Client, dbName string, date string) error {
	// Delete the person of the day from the Persons collection
	collection := client.Database(dbName).Collection("GuessesOfTheMonth")
	if collection == nil {
		return fmt.Errorf("persons collection not found")
	}
	// Format the date as YYYY-MM-DD
	dateFrom := date // Date from one day ago

	// Delete the person of the day for the current date
	_, err := collection.DeleteOne(context.TODO(), map[string]interface{}{"date": dateFrom})
	if err != nil {
		return fmt.Errorf("failed to delete person of the day: %v", err)
	}

	return nil
}

func InitDB(mongoClient *mongo.Client) error {
	CreateDatabase(mongoClient, "dodle")
	fmt.Println("Creating collections...")
	CreateCollection(mongoClient, "dodle", "Persons")
	CreateCollection(mongoClient, "dodle", "GuessesOfTheMonth")

	// Load persons from file
	persons, err := data.OpenPersonsFile()
	if err != nil {
		return fmt.Errorf("failed to open persons file: %v", err)
	}

	// Check if collection is empty before populating
	existingPersons, err := GetPersons(mongoClient, "dodle")
	if err != nil {
		return fmt.Errorf("failed to check existing persons: %v", err)
	}

	if len(existingPersons.Persons) == 0 {
		fmt.Println("Populating Persons collection...")
		result := PopulatePersonsCollection(mongoClient, "dodle", persons)
		if result != "" {
			return fmt.Errorf("failed to populate persons collection: %s", result)
		}
	} else {
		fmt.Println("Persons collection already contains data, skipping population")
	}

	return nil
}

func UpdatePersonOfTheDay(mongoClient *mongo.Client) error {
	previousPersons, err := GetPersonsOfTheDay(mongoClient, "dodle")
	if err != nil {
		return fmt.Errorf("failed to get previous persons of the day: %v", err)
	}

	personsAvailable, err := GetPersons(mongoClient, "dodle")
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
	if GetPersonOfTheDay, err := GetPersonOfTheDay(mongoClient, "dodle"); err == nil {
		if GetPersonOfTheDay.Firstname != "" {
			if err := DeletePersonOfTheDay(mongoClient, "dodle", dateOfToday); err != nil {
				return fmt.Errorf("failed to delete previous person of the day: %v", err)
			}
			fmt.Println("Previous person of the day deleted successfully")
		}
	}

	if err := CreatePersonOfTheDay(mongoClient, "dodle", candidate); err != nil {
		return fmt.Errorf("failed to create person of the day: %v", err)
	}

	dateToDelete := time.Now().AddDate(0, 0, -10).Format("2006-01-02")

	if err := DeletePersonOfTheDay(mongoClient, "dodle", dateToDelete); err != nil {
		return fmt.Errorf("failed to delete previous person of the day: %v", err)
	}
	fmt.Println("Previous person of the day deleted successfully")
	// You can add your logic here to update the person of the day in the database.
	return nil
}

func GetPersonOfYesterday(client *mongo.Client, dbName string) (persons.Person, error) {
	// Retrieve the person of yesterday from the GuessesOfTheMonth collection
	collection := client.Database(dbName).Collection("GuessesOfTheMonth")
	if collection == nil {
		return persons.Person{}, fmt.Errorf("GuessesOfTheMonth collection not found")
	}

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02") // Format the date as YYYY-MM-DD

	var doc struct {
		Date   string         `bson:"date"`
		Person persons.Person `bson:"person"`
	}

	err := collection.FindOne(context.TODO(), map[string]interface{}{"date": yesterday}).Decode(&doc)
	if err != nil {
		return persons.Person{}, fmt.Errorf("failed to find person of yesterday: %v", err)
	}

	return doc.Person, nil
}

func GetGuessID(client *mongo.Client, dbName string) (string, error) {
	// Retrieve the ID of the guess from the GuessesOfTheMonth collection
	collection := client.Database(dbName).Collection("GuessesOfTheMonth")
	if collection == nil {
		return "", fmt.Errorf("GuessesOfTheMonth collection not found")
	}

	cursor, err := collection.Find(context.TODO(), map[string]interface{}{})
	if err != nil {
		return "", fmt.Errorf("failed to find persons of the day: %v", err)
	}
	defer cursor.Close(context.TODO())

	var doc struct {
		ID     string         `bson:"_id"`
		Date   string         `bson:"date"`
		Person persons.Person `bson:"person"`
	}

	var guesses []string
	for cursor.Next(context.TODO()) {

		if err := cursor.Decode(&doc); err != nil {
			return "", fmt.Errorf("failed to decode document: %v", err)
		}

		guesses = append(guesses, doc.ID)
	}

	if (len(guesses) == 0) {
		return "", fmt.Errorf("no guesses found for today")
	}
	
	return guesses[len(guesses)-1], nil
}
