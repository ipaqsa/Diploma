package kernel

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

//go get github.com/mattn/go-sqlite3

type DB struct {
	ptr *sql.DB
	mtx sync.Mutex
}

func (client *Client) DBDialogsInit() {
	db, err := sql.Open("sqlite3", "dialogs.db")
	if err != nil {
		return
	}

	client.dbDialogs = &DB{
		ptr: db,
	}
}

//func (client *Client) CreateDialogTable(dialog string)  {
//	_, err := client.dbDialogs.ptr.Exec(
//		`CREATE TABLE IF NOT EXISTS $1 (
//    	 VARCHAR(75) UNIQUE,
//    	data VARCHAR(500),
//    	password VARCHAR(25),
//    	name VARCHAR(30),
//    	room VARCHAR(1),
//    	PRIMARY KEY(login)
//    	);
//	`, dialog)
//	if err != nil {
//		return
//	}
//}

func (client *Client) DBUsersInit() {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		println("ERROR: Create DB")
		return
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS users (
    	login VARCHAR(75) UNIQUE,
    	key VARCHAR(500),
    	password VARCHAR(25),
    	name VARCHAR(30),
    	room VARCHAR(1),
    	PRIMARY KEY(login)
    	);
	`)
	if err != nil {
		println(err.Error())
		return
	}
	client.dbUsers = &DB{
		ptr: db,
	}
}

func DBFriendsInit(filename string) *DB {
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

func (client *Client) SaveUser(user *User) error {
	if client.dbUsers.ptr == nil {
		return errors.New("ERROR: BD ptr is nil")
	}
	client.dbUsers.mtx.Lock()
	defer client.dbUsers.mtx.Unlock()
	_, err := client.dbUsers.ptr.Exec(`INSERT INTO users (login, key, password, name, room) VALUES ($1, $2, $3, $5, $6)`,
		user.Login, StringPrivate(user.PrivateKey), user.Password, user.Name, user.Room)
	return err
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
func GetUserFromDB(user *User, password string) uint {
	psswd := HashSum([]byte(password))
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return 0
	}
	var string_key string
	row := db.QueryRow(`SELECT * FROM users WHERE login=$1 AND password=$2 LIMIT 1`, user.Login, psswd)
	err = row.Scan(&user.Login, &string_key, &user.Password, &user.Name, &user.Room)
	if err != nil {
		return 0
	}
	user.PrivateKey = ParsePrivate(string_key)
	return 1
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
	client.dbFriends.mtx.Lock()
	defer client.dbFriends.mtx.Unlock()
	var members []string
	rows, err := client.dbFriends.ptr.Query(`SELECT login FROM friendsLogins`)
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
