package gopeer

import (
	"crypto/rsa"
	"net"
	"sync"
)

type Listener struct {
	client *Client
	listen net.Listener
}

type User struct {
	Name       string          `json:"name"`
	Login      string          `json:"login"`
	Password   []byte          `json:"password"`
	Room       uint            `json:"room"`
	PrivateKey *rsa.PrivateKey `json:"privateKey"`
}

type Client struct {
	dbFriends   *DB
	dbUsers     *DB
	user        *User
	address     string
	mapping     map[string]bool
	connections map[net.Conn]string
	actions     map[string]chan string
	mutex       *sync.Mutex
	f2f         map[string]*rsa.PublicKey
	f2f_d       map[*rsa.PublicKey]string
}

type Package struct {
	Head HeadPackage `json:"head"`
	Body BodyPackage `json:"body"`
}

type HeadPackage struct {
	Rand    string `json:"rand"`
	Title   string `json:"title"`
	Sender  string `json:"sender"`
	Session string `json:"session"`
}

type BodyPackage struct {
	Date string `json:"date"`
	Data string `json:"data"`
	Hash string `json:"hash"`
	Sign string `json:"sign"`
}

type UserBroadcast struct {
	Name       string                    `json:"name"`
	Login      string                    `json:"login"`
	Password   []byte                    `json:"password"`
	Room       uint                      `json:"room"`
	PrivateKey *rsa.PrivateKey           `json:"privateKey"`
	F2F        map[string]*rsa.PublicKey `json:"f2f"`
}

type PackageBroadcast struct {
	Login            string
	Address          string
	AddressBroadcast string
	Key              string
	Room             uint
}

type NodeScanner struct {
	login       string
	db          *DB
	Key         string
	Room        uint
	Connections map[string]string
	Address     string
	AddressR    string
}

const (
	END_BYTES = "\000\001\002\003\004\005"
)
