package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func InitDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}

	err := db.init()

	return db, err

}

func (db *DB) WriteDB(structure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := json.Marshal(structure)

	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0600)

	if err != nil {
		return err
	}

	return nil

}

func (db *DB) CreateChirp(body string) (Chirp, error) {

	dbStructure, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	nextId := len(dbStructure.Chirps) + 1

	chirp := Chirp{
		Id:   nextId,
		Body: body,
	}

	dbStructure.Chirps[nextId] = chirp

	err = db.WriteDB(dbStructure)

	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil

}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()

	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))

	for _, item := range dbStructure.Chirps {
		chirps = append(chirps, item)
	}

	return chirps, nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := DBStructure{}

	dbData, err := os.ReadFile(db.path)

	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}

	err = json.Unmarshal(dbData, &dbStructure)

	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil

}

func (db *DB) init() error {
	_, err := os.ReadFile(db.path)

	if errors.Is(err, os.ErrNotExist) {
		return db.createDBFile()
	}

	return err
}

func (db *DB) createDBFile() error {
	structure := DBStructure{
		Chirps: map[int]Chirp{},
	}

	return db.WriteDB(structure)
}
