package kernel

import (
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"sync"
)

//go get github.com/mattn/go-sqlite3

type DB struct {
	ptr *sql.DB
	mtx sync.Mutex
}

func (client *Client) DBDialogsInit() {
	db, err := sql.Open("sqlite3", "./data/dialogs.db")
	if err != nil {
		return
	}
	client.dbDialogs = &DB{
		ptr: db,
	}
}

func (client *Client) CreateDialogTable(dialog string) {
	client.dbDialogs.mtx.Lock()
	defer client.dbDialogs.mtx.Unlock()
	query := "CREATE TABLE IF NOT EXISTS " + dialog + " (id INTEGER PRIMARY KEY AUTOINCREMENT, date VARCHAR(25), sender VARCHAR(30), data TEXT);"
	_, err := client.dbDialogs.ptr.Exec(query)
	if err != nil {
		println(err.Error())
		return
	}
}

func (client *Client) AddMessage(dialog string, pack *Package) {
	client.dbDialogs.mtx.Lock()
	defer client.dbDialogs.mtx.Unlock()
	message := &Message{
		Date: pack.Body.Date,
		Data: pack.Body.Data,
	}
	if pack.Head.Sender == "" {
		message.From = client.user.Login
	} else {
		message.From = client.GetLogin(pack.Head.Sender)
	}
	query := "INSERT INTO " + dialog + " (date, sender, data) VALUES ($1, $2, $3)"

	_, err := client.dbDialogs.ptr.Exec(query, message.Date, message.From, message.Data)
	if err != nil {
		return
	}
}

func (client *Client) DBUsersInit() {
	client.f2f[client.user.Login] = client.Public()
	db, err := sql.Open("sqlite3", "./data/users.db")
	if err != nil {
		println(err.Error())
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

func (client *Client) DBHashesInit() {
	db, err := sql.Open("sqlite3", "./data/externals.db")
	if err != nil {
		return
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS hashes (
    	login VARCHAR(75) UNIQUE,
    	date VARCHAR(25),
    	hash VARCHAR(500),
    	PRIMARY KEY(login)
    	);
	`)
	if err != nil {
		return
	}
	client.dbExternals = &DB{
		ptr: db,
	}
}

func (client *Client) AddHash(msg, date, login string) error {
	var user User
	err := json.Unmarshal(Base64Decode(msg), &user)
	println(user.Name)
	if err != nil {
		return err
	}
	sumUser, err := HashSumUser(&user)
	if err != nil {
		return err
	}
	client.dbExternals.mtx.Lock()
	defer client.dbExternals.mtx.Unlock()
	_, err = client.dbExternals.ptr.Exec(`INSERT INTO hashes (login, date, hash) VALUES ($1, $2, $3)`, login, date, sumUser)
	return err
}

func (client *Client) DBFriendsInit() {
	if exists("./data/friends.db") == false {
		os.Remove("friends.db")
	}
	db, err := sql.Open("sqlite3", "./data/friends.db")
	if err != nil {
		return
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS friendsLogins (
    	login VARCHAR(75) UNIQUE,
    	key VARCHAR(500),
    	PRIMARY KEY(login)
    	);
	`)
	if err != nil {
		return
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS friendsAddresses (
    	key VARCHAR(500) UNIQUE,
    	address VARCHAR(25),
    	PRIMARY KEY(key)
    	);
	`)
	if err != nil {
		return
	}
	client.dbFriends = &DB{
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
	psswd := Base64Encode(HashSum([]byte(password)))
	var stringKey string
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		println(err.Error())
		return 0
	}
	row := db.QueryRow(`SELECT * FROM users WHERE login=$1 AND password=$2 LIMIT 1`, user.Login, psswd)
	err = row.Scan(&user.Login, &stringKey, &user.Password, &user.Name, &user.Room)
	if err != nil {
		return 0
	}
	user.PrivateKey = ParsePrivate(stringKey)
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

func (client *Client) GetDialog(dialogName string) []Message {
	client.dbDialogs.mtx.Lock()
	defer client.dbDialogs.mtx.Unlock()
	var dialog []Message
	query := "SELECT * FROM " + dialogName
	rows, err := client.dbDialogs.ptr.Query(query)
	if err != nil {
		return nil
	}
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.Date, &message.From, &message.Data); err != nil {
			return dialog
		}
		dialog = append(dialog, message)
	}
	return dialog
}

func (client *Client) GetHashFromDialog(dialogName string) string {
	var summary string
	dialog := client.GetDialog(dialogName)
	for _, msg := range dialog {
		msgSummary := msg.Data + msg.Date + msg.From
		summary += msgSummary
	}
	return Base64Encode(HashSum([]byte(summary)))
}

func HashSumUser(user *User) (string, error) {
	marshal, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	return Base64Encode(HashSum(marshal)), err
}
