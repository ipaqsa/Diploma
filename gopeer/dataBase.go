package gopeer

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

//go get github.com/mattn/go-sqlite3

type DB struct {
	ptr *sql.DB
	mtx sync.Mutex
}

func DBInit(filename string) *DB {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS friendsLogins (
    	login VARCHAR(75) UNIQUE,
    	key VARCHAR(500),
    	PRIMARY KEY(login)
    	);
	`)
	if err != nil {
		return nil
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS friendsAddresses (
    	key VARCHAR(500) UNIQUE,
    	address VARCHAR(25),
    	PRIMARY KEY(key)
    	);
	`)
	if err != nil {
		return nil
	}
	return &DB{
		ptr: db,
	}
}

func (db *DB) SetLogin(login, key string) error {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	_, err := db.ptr.Exec(`INSERT INTO friendsLogins (login, key) VALUES ($1, $2)`, login, key)
	return err
}

func (db *DB) SetAddress(key, address string) error {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	_, err := db.ptr.Exec(`INSERT INTO friendsAddresses (key, address) VALUES ($1, $2)`, key, address)
	return err
}

func (db *DB) GetKey(login string) string {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	var key string
	row := db.ptr.QueryRow(`SELECT key FROM friendsLogins WHERE login=$1 LIMIT 1`, login)
	err := row.Scan(&key)
	if err != nil {
		return ""
	}
	return key
}

func (db *DB) GetAddress(key string) string {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	var address string
	row := db.ptr.QueryRow(`SELECT address FROM friendsAddresses WHERE key=$1 LIMIT 1`, key)
	err := row.Scan(&address)
	if err != nil {
		return ""
	}
	return address
}

func (db *DB) SizeLogins() int {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	var data int
	row := db.ptr.QueryRow(`SELECT COUNT(*) FROM friendsLogins`)
	err := row.Scan(&data)
	if err != nil {
		return 0
	}
	return data
}

func (client *Client) GetAllMembers() []string {
	client.db.mtx.Lock()
	defer client.db.mtx.Unlock()
	var members []string
	rows, err := client.db.ptr.Query(`SELECT login FROM friendsLogins`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var member string
		if err := rows.Scan(&member); err != nil {
			return members
		}
		members = append(members, member)
	}
	return members
}
