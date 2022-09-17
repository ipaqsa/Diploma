package gopeer

import (
	"encoding/json"
	"net"
	"time"
)

var infoLoggerBroadcast = newLogger("broadcast", "INFO")
var errorLoggerBroadcast = newLogger("broadcast", "ERROR")

func (client *Client) NewNodeBroadcast(address, addressR, login, key string, room uint) *NodeScanner {
	dbname := login + "Friends" + ".db"
	db := DBInit(dbname)
	client.db = db
	infoLoggerBroadcast.Printf("Node was created")
	return &NodeScanner{
		login:       login,
		db:          db,
		Key:         key,
		Room:        room,
		Connections: make(map[string]string),
		Address:     address,
		AddressR:    addressR,
	}
}

func (node *NodeScanner) Run(addresses []string) {
	go handleBroadcastServer(node)
	time.Sleep(time.Second * 2)
	for {
		if addresses != nil {
			for _, address := range addresses {
				if address != node.Address {
					node.SendToAddress(address)
				}
			}
		}
		time.Sleep(time.Second * 15)
	}
}

func handleBroadcastServer(node *NodeScanner) {
	listen, err := net.Listen("tcp", node.AddressR)
	if err != nil {
		errorLoggerBroadcast.Printf("Server creating error")
		panic("Server Error")
	}
	defer listen.Close()
	infoLoggerBroadcast.Printf("BroadcastServer was started with ", node.AddressR)
	for {
		conn, err := listen.Accept()
		infoLoggerBroadcast.Printf("New connection")
		if err != nil {
			break
		}
		go handleConnection(node, conn)
	}
}

func handleConnection(node *NodeScanner, conn net.Conn) {
	defer conn.Close()
	var (
		buffer  = make([]byte, 512)
		message string
		pack    PackageBroadcast
	)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			break
		}
		message += string(buffer[:length])
	}
	err := json.Unmarshal([]byte(message), &pack)
	if err != nil {
		return
	}
	if node.Room == pack.Room && node.login != pack.Login {
		if node.db.GetKey(pack.Login) == "" {
			//println("Save", pack.Login)
			node.db.SetLogin(pack.Login, pack.Key)
			node.db.SetAddress(pack.Key, pack.Address)
			node.SendToAddress(pack.Address)
		} else if node.db.GetKey(pack.Login) != pack.Key {
			//println("Update data", pack.Login)
			node.db.SetLogin(pack.Login, pack.Key)
			node.db.SetAddress(pack.Key, pack.Address)
		} else if node.db.GetKey(pack.Login) != pack.Address {
			//println("Update address", pack.Login)
			node.db.SetAddress(pack.Key, pack.Address)
		}
	}
}

func (node *NodeScanner) SendToAddress(address string) {
	var pack = PackageBroadcast{
		Login:   node.login,
		Address: node.Address,
		Key:     node.Key,
		Room:    node.Room,
	}
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return
	}
	defer conn.Close()
	js_pack, _ := json.Marshal(pack)
	conn.Write(js_pack)
	infoLogger.Printf("Message was sent to", address)
}
