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
		print(err)
		return &DB{}, err
	}
	return &db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, id int) (Chirp, error) {
	err := db.ensureDB()
	if err != nil {
		print("Ensure DB failed in CreateChirp")
		return Chirp{}, err
	}
	dbs, err := db.loadDB()
	if err != nil {
		print("Load db failed in create chirp")
		return Chirp{}, err
	}
	dbs.Chirps[id] = Chirp{Body: body}
	db.writeDB(dbs)
	return dbs.Chirps[id], nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	err := db.ensureDB()
	if err != nil {
		print("EnsureDB failed in get chirps")
		return []Chirp{}, err
	}
	dbs, err := db.loadDB()
	if err != nil {
		print("load db failed in get chirps")
		return []Chirp{}, err
	}
	chirps := make([]Chirp, len(dbs.Chirps))
	for _, v := range dbs.Chirps {
		chirps = append(chirps, v)
	}
	return chirps, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	_, err := os.Create(db.path)
	if err != nil {
		print(err)
		return err
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	err := db.ensureDB()
	if err != nil {
		print(err)
		return DBStructure{}, err
	}
	read, err := os.ReadFile(db.path)
	if err != nil {
		print(err)
		return DBStructure{}, err
	}

	var dbStructure DBStructure
	if len(read) == 0 {
		return DBStructure{
			Chirps: map[int]Chirp{},
		}, nil
	} else {
		err = json.Unmarshal(read, &dbStructure.Chirps)
		if err != nil {
			print(err)
			return DBStructure{}, err
		}
	}
	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	err := db.ensureDB()
	if err != nil {
		print(err)
		return err
	}
	write, err := json.Marshal(dbStructure.Chirps)
	if err != nil {
		print(err)
		return err
	}
	os.WriteFile(db.path, write, 0666)
	return nil
}
