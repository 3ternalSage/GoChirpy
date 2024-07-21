package database

import (
	"encoding/json"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Body string `json:"body"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	var db DB
	db.path = path
	db.mux = &sync.RWMutex{}
	err := db.ensureDB()
	if err != nil {
		return &DB{}, err
	}
	return &db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	ch := Chirp{
		Body: body,
	}
	json, err := json.Marshal(ch)
	if err != nil {
		return Chirp{}, err
	}
	os.WriteFile(db.path, []byte(json), 0)
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error)

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	err := os.WriteFile(db.path, []byte{}, os.FileMode(0))
	if err != nil {
		return err
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error)

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error
