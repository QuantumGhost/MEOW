package main

import (
    "io/ioutil"
    "testing"
)

func TestSQLiteStorage(t *testing.T) {
	storage := NewSQLiteStorage(":memory:")
	conn := storage.conn
	
	// migrate database.
	dat, err := ioutil.ReadFile("database/auth.sql")
	if err != nil {
		t.Error("Error when loading schema sql file.", err)
	}
	
	_, err = conn.Exec(string(dat))	
	if err != nil {
		t.Error("Error while creating table.")
	}
	
	var testData = []struct {
		username string
		password string 
		port uint16
		active uint8 		
	}{
		{"alice", "alice", 0, 1},
		{"bob", "bob", 1024, 1},
		{"david", "david", 0, 0},
		{"eve", "eve", 1024, 0},
	}
	
	const INSERT_SQL = "INSERT INTO auth (username, password, port, active) VALUES (?, ?, ?, ?)"
	stmt, err := conn.Prepare(INSERT_SQL)
	if err != nil {
		t.Error("Error while prepare an insert statement", err)
	}
	for _, td := range testData {
		_, err := stmt.Exec(td.username, td.password, td.port, td.active)
		if err != nil {
			t.Errorf("Error while inserting %s, error: %s", td, err)
		}
	}
	
	// run tests
	for _, td := range testData{
		au, ok := storage.Get(td.username)
		if td.active == 1 && ok != nil {
			t.Errorf("Get should return a record for %s.", td.username)
		}
		if td.active == 0 && ok == nil {
			t.Errorf("Get should return nil for %s.", td.username)
		}
		if au != nil {
			if au.passwd != td.password {
				t.Errorf("Password is not the same for %s.", td.username)
			}
			if au.port != td.port {
				t.Errorf("Port is not the same for %s.", td.username)
			}
		}
	}
}

func TestMemoryStorage(t *testing.T) {
    storage := NewMemoryStorage()
    
    var testData = []struct {
		username string
		password string 
		port uint16
		active uint8 		
	}{
		{"alice", "alice", 0, 1},
		{"bob", "bob", 1024, 1},
		{"david", "david", 0, 0},
		{"eve", "eve", 1024, 0},
	}
    
    for _, td := range testData {
        storage.Set(td.username, td.password, td.port, td.active != 0)
        au, ok := storage.Get(td.username)
        if ok != nil {
            t.Errorf("Get should return a record for %s", td.username)
        }
        
        if au.passwd != td.password || au.port != td.port {
            t.Error("Get should return a authUser with same values.")
        }
        storage.Delete(td.username)
        _, ok = storage.Get(td.username)
        if ok == nil {
            t.Error("Get a deleted authUser should return userNotExist.") 
        }
    }
    
    td := testData[0]
    storage.Set(td.username, td.password, td.port, td.active != 0)
    storage.Deactivate(td.username)
    _, ok := storage.Get(td.username)
    if ok == nil {
        t.Error("Get a deactivated authUser should return userNotExist.")
    }
    
    _, ok = storage.Get("henry")
    if ok == nil {
        t.Error("Get a non-existing authUser should return userNotExist.")
    }
}