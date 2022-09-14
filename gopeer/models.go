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
	Name       string                    `json:"name"`
	Login      string                    `json:"login"`
	Password   []byte                    `json:"password"`
	Room       uint                      `json:"room"`
	PrivateKey *rsa.PrivateKey           `json:"privateKey"`
	F2F        map[string]*rsa.PublicKey `json:"f2f"`
}

type Client struct {
	user        *User
	address     string
	mapping     map[string]bool
	connections map[net.Conn]string
	actions     map[string]chan string
	mutex       *sync.Mutex
	f2f_d       map[*rsa.PublicKey]string
}

type Package struct {
	Head HeadPackage `json:"head"`
	Body BodyPackage `json:"body"`
}

type PackageBroadCast struct {
	AddressFrom string
	Key         *rsa.PublicKey
	Room        uint
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
