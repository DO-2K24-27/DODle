package data

import (
	persons "api/struct"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func OpenPersonsFile() (persons.Persons, error) {
	jsonFile, err := os.Open("./data/persons.json")
	if err != nil {
		return persons.Persons{}, fmt.Errorf("error opening file: %v", err)
	}
	defer func() {
		if err := jsonFile.Close(); err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}()

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
