package database

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PublicUser struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
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

func (db *DB) GetChirpById(id int) (Chirp, error) {
	structure, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	myChirp, ok := structure.Chirps[id]

	if !ok {
		return myChirp, errors.New("there Is No Chirp with this id")
	}

	return myChirp, nil

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
		Users:  map[int]User{},
	}

	return db.WriteDB(structure)
}

// users

func (db *DB) CreateUser(email string, password string) (PublicUser, error) {

	dbStructure, err := db.loadDB()

	if err != nil {
		return PublicUser{}, err
	}
	for _, v := range dbStructure.Users {
		if v.Email == email {
			return PublicUser{}, errors.New("user with this email already exist")
		}
	}
	hashBytesPass, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return PublicUser{}, err
	}
	hashedPassword := string(hashBytesPass)

	nextId := len(dbStructure.Users) + 1

	user := User{
		Id:       nextId,
		Email:    email,
		Password: hashedPassword,
	}

	dbStructure.Users[nextId] = user

	err = db.WriteDB(dbStructure)

	if err != nil {
		return PublicUser{}, err
	}

	return PublicUser{
		Id:    user.Id,
		Email: user.Email,
	}, nil

}

func (db *DB) UpdateUser(email string, password string, id string) (PublicUser, error) {

	dbStructure, err := db.loadDB()

	if err != nil {
		return PublicUser{}, err
	}

	err = db.WriteDB(dbStructure)

	if err != nil {
		return PublicUser{}, err
	}

	idNum, err := strconv.Atoi(id)

	if err != nil {
		return PublicUser{}, err
	}

	if dbStructure.Users[idNum].Id != idNum {
		return PublicUser{}, errors.New("user with this id doesn't exist")
	}

	User := User{
		Id:       idNum,
		Email:    email,
		Password: password,
	}

	dbStructure.Users[idNum] = User

	err = db.WriteDB(dbStructure)

	if err != nil {
		return PublicUser{}, err
	}

	return PublicUser{
		Id:    idNum,
		Email: email,
	}, nil

}

func (db *DB) LoginVerification(email string, password string) (PublicUser, error) {

	dbStructure, err := db.loadDB()

	user := User{}

	if err != nil {
		return PublicUser{}, err
	}

	for _, usrV := range dbStructure.Users {
		if usrV.Email == email {
			user = usrV
			break
		}

	}

	if user.Email == "" {
		return PublicUser{}, errors.New("wrong login credentials entered")
	}

	// compare password

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return PublicUser{}, errors.New("wrong login credentials entered")
	}

	return PublicUser{
		Id:    user.Id,
		Email: user.Email,
	}, nil

}
