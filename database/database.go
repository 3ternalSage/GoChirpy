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
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbs, err := db.LoadDB()
	if err != nil {
		print("Load db failed in create chirp")
		return Chirp{}, err
	}
	id := len(dbs.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}
	dbs.Chirps[id] = chirp

	err = db.writeDB(dbs)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbs, err := db.LoadDB()
	if err != nil {
		print("load db failed in get chirps")
		return []Chirp{}, err
	}
	chirps := make([]Chirp, 0, len(dbs.Chirps))
	for _, v := range dbs.Chirps {
		chirps = append(chirps, v)
	}
	return chirps, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if err != nil {
		return db.writeDB(DBStructure{
			Chirps: make(map[int]Chirp),
		})
	}
	return nil
}

// LoadDB reads the database file into memory
func (db *DB) LoadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure := DBStructure{}
	print("Reading file from path: ", db.path, "\n")
	read, err := os.ReadFile(db.path)
	if err != nil {
		print("Read file failed in load db\n")
		print(err)
		return dbStructure, err
	}
	err = json.Unmarshal(read, &dbStructure)
	if err != nil {
		print("unmarshal failed in db\n")
		print(err)
		return DBStructure{}, err
	}
	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	write, err := json.Marshal(dbStructure)
	if err != nil {
		print(err)
		return err
	}
	err = os.WriteFile(db.path, write, 0600)
	if err != nil {
		print(err)
		return err
	}
	return nil
}
