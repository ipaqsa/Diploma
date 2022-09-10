package gopeer

import (
	"crypto/rsa"
	"net"
	"sync"
)

type Listener struct {
	address string
	client  *Client
	listen  net.Listener
}

type Client struct {
	privateKey  *rsa.PrivateKey
	mapping     map[string]bool
	connections map[net.Conn]string
	actions     map[string]chan string
	mutex       *sync.Mutex
	f2f         FriendToFriend
}

type FriendToFriend struct {
	enable  bool
	friends map[string]*rsa.PublicKey
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
	Data string `json:"data"`
	Hash string `json:"hash"`
	Sign string `json:"sign"`
}
