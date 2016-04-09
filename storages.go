package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	
)

var userNotExist = errors.New("User does not found")

// PasswordStorage stores all secret needed in authentication.
// A basic Interface for password storage.
type AuthStorage interface {
	setup()   // 初始化函数
	Get(username string) (*authUser, error)
	Set(username string, password string, port uint16, active bool) error
	Delete(username string) error
	Deactivate(username string) error
}

// MemoryStorage type
// A simple storage which stores all password in a map. 
type MemoryStorage struct {
	credentials map[string]*authUser
}

func NewMemoryStorage() *MemoryStorage {
	storage := new(MemoryStorage)
	credentials := make(map[string]*authUser)
	storage.credentials = credentials
	return storage
}

func (storage *MemoryStorage) setup() {}

// Get MemoryStorage type
// 
func (storage *MemoryStorage) Get(username string) (*authUser, error) {
	au, ok := storage.credentials[username] 
	if !ok {
		return nil, userNotExist
	}
	return au, nil
}

func (storage *MemoryStorage) Set(
		username string, password string, port uint16, active bool) error {
	au := &authUser{password, "", port}
	storage.credentials[username] = au
	return nil
}

func (storage *MemoryStorage) Delete(username string) error {
	delete(storage.credentials, username)
	return nil
}

func (storage *MemoryStorage) Deactivate(username string) error {
	delete(storage.credentials, username)
	return nil
}

// SQLiteStorage type
// store all password in a given database.
type SQLiteStorage struct {
	file string
	conn *sql.DB
}

func NewSQLiteStorage(file string) *SQLiteStorage {
	storage := new(SQLiteStorage)
	storage.file = file
	storage.setup()
	return storage
}

func (self *SQLiteStorage) setup() {
	db, err := sql.Open("sqlite3", self.file)
	if err != nil {
		Fatal("Error while opening database file", err)
	}
	self.conn = db
}

func (self *SQLiteStorage) Get(username string) (*authUser, error) {
	const SELECT_SQL = `
	SELECT password, port FROM auth 
	WHERE username=? AND active = 1 LIMIT 1
	`

	var password string 
	var port uint16
	
	err := self.conn.QueryRow(SELECT_SQL, username).Scan(&password, &port)
	switch {
		case err == sql.ErrNoRows:
			debug.Printf("Not exists")
			return nil, userNotExist
		case err != nil:
			Fatal("Error while getting user from database:\n", err)
	}
	au := &authUser{password, "", port} 
	return au, nil
}

func (self *SQLiteStorage) Set(
		username string, password string, port uint16, active bool) error {
	return nil
}

func (self *SQLiteStorage) Delete(username string) error {
	return nil
}

func (self *SQLiteStorage) Deactivate(username string) error {
	return nil
}
