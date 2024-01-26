package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type RevokedAccessToken struct {
	Id        string    `json:"id"`
	RevokedAt time.Time `json:"revokedAt"`
}

type PublicUser struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type DBStructure struct {
	Chirps              map[int]Chirp                 `json:"chirps"`
	Users               map[int]User                  `json:"users"`
	RevokedAccessTokens map[string]RevokedAccessToken `json:"revokedAccessTokens"`
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

func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {

	dbStructure, err := db.loadDB()

	if err != nil {
		return Chirp{}, err
	}

	nextId := len(dbStructure.Chirps) + 1

	chirp := Chirp{
		Id:       nextId,
		Body:     body,
		AuthorId: authorId,
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

// tokens

func (db *DB) CheckIsTokenRevoked(tokenId string) (bool, error) {

	structure, err := db.loadDB()

	if err != nil {
		return true, err
	}

	_, ok := structure.RevokedAccessTokens[tokenId]

	if ok {
		return true, nil
	}

	return false, nil

}

func (db *DB) RevokeAccessToken(tokenId string) error {

	structure, err := db.loadDB()

	if err != nil {
		return err
	}

	_, ok := structure.RevokedAccessTokens[tokenId]

	if ok {
		return errors.New("token already revoked")
	}

	token := RevokedAccessToken{
		Id:        tokenId,
		RevokedAt: time.Now().UTC(),
	}

	structure.RevokedAccessTokens[tokenId] = token

	err = db.WriteDB(structure)

	if err != nil {
		return err
	}

	return nil

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
		Chirps:              map[int]Chirp{},
		Users:               map[int]User{},
		RevokedAccessTokens: map[string]RevokedAccessToken{},
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
		Id:          nextId,
		Email:       email,
		Password:    hashedPassword,
		IsChirpyRed: false,
	}

	dbStructure.Users[nextId] = user

	err = db.WriteDB(dbStructure)

	if err != nil {
		return PublicUser{}, err
	}

	return PublicUser{
		Id:          user.Id,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}, nil

}

func (db *DB) UpdateUser(email string, password string, id int, isChirpyRed *bool) (PublicUser, error) {

	dbStructure, err := db.loadDB()

	if err != nil {
		return PublicUser{}, err
	}

	err = db.WriteDB(dbStructure)

	if err != nil {
		return PublicUser{}, err
	}

	if dbStructure.Users[id].Id != id {
		return PublicUser{}, errors.New("user with this id doesn't exist")
	}

	if email == "" {
		email = dbStructure.Users[id].Email
	}

	if password == "" {
		password = dbStructure.Users[id].Password
	}

	isRed := dbStructure.Users[id].IsChirpyRed

	if isChirpyRed != nil {
		isRed = *isChirpyRed
	}

	user := User{
		Id:          id,
		Email:       email,
		Password:    password,
		IsChirpyRed: isRed,
	}

	dbStructure.Users[id] = user

	err = db.WriteDB(dbStructure)

	if err != nil {
		return PublicUser{}, err
	}

	return PublicUser{
		Id:    id,
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
		Id:          user.Id,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}, nil

}

func (db *DB) DeleteChirp(chirpID int, userId int) error {

	dbStructure, err := db.loadDB()

	if err != nil {
		return err
	}

	if dbStructure.Chirps[chirpID].AuthorId != userId {
		return errors.New("you can only delete your own chirps")
	}

	delete(dbStructure.Chirps, chirpID)

	err = db.WriteDB(dbStructure)

	if err != nil {
		return err
	}

	return nil

}
