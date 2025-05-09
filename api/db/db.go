package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

type Person struct {
	Data json.RawMessage `json:"data"`
}

type DB struct {
	*sql.DB
}

func NewDB(connStr string) (*DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return &DB{db}, nil
}

func (db *DB) Init() error {
	// Read and execute migration file
	migrationPath := filepath.Join("db", "migrations", "001_create_persons_table.sql")
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("error reading migration file: %v", err)
	}

	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		return fmt.Errorf("error executing migration: %v", err)
	}

	return nil
}

func (db *DB) SeedData() error {
	// Read persons.json file
	data, err := os.ReadFile("./data/persons.json")
	if err != nil {
		return fmt.Errorf("error reading persons.json: %v", err)
	}

	// Parse JSON array
	var persons []Person
	if err := json.Unmarshal(data, &persons); err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}
	defer tx.Rollback()

	// Clear existing data
	_, err = tx.Exec("TRUNCATE TABLE persons")
	if err != nil {
		return fmt.Errorf("error truncating table: %v", err)
	}

	// Insert new data
	stmt, err := tx.Prepare("INSERT INTO persons (data) VALUES ($1)")
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	for _, person := range persons {
		_, err = stmt.Exec(person.Data)
		if err != nil {
			return fmt.Errorf("error inserting data: %v", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
} 

func (db *DB) GetPersons() ([]Person, error) {
	rows, err := db.Query("SELECT data FROM persons")
	if err != nil {
		return nil, fmt.Errorf("error querying persons: %v", err)
	}

	defer rows.Close()

	var persons []Person
	for rows.Next() {
		var person Person
		if err := rows.Scan(&person.Data); err != nil {
			return nil, fmt.Errorf("error scanning person: %v", err)
		}
		persons = append(persons, person)
	}

	return persons, nil
}
