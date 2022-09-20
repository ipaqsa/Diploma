package kernel

import (
	"encoding/json"
	"net"
	"strings"
	"time"
)

var infoLoggerBroadcast = newLogger("broadcast", "INFO")
var errorLoggerBroadcast = newLogger("broadcast", "ERROR")

func (client *Client) NewNodeBroadcast(address, addressR, login, key string, room uint) *NodeScanner {
	dbname := login + "Friends" + ".db"
	db := DBFriendsInit(dbname)
	client.dbFriends = db
	infoLoggerBroadcast.Printf("Node was created")
	return &NodeScanner{
		login:       login,
		db:          db,
		Key:         key,
		Room:        room,
		Connections: make(map[string]string),
		Address:     address,
		AddressB:    addressR,
	}
}

func (node *NodeScanner) BroadcastMSG() {
	conn, err := net.ListenPacket("udp4", ":9001")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	addr, err := net.ResolveUDPAddr("udp4", "192.168.0.255:9001")
	if err != nil {
		errorLoggerBroadcast.Printf(err.Error())
	}
	var pack = PackageBroadcast{
		Login:            node.login,
		Address:          node.Address,
		AddressBroadcast: node.AddressB,
		Key:              node.Key,
		Room:             node.Room,
	}

	JSONPack, _ := json.Marshal(pack)
	_, err = conn.WriteTo(JSONPack, addr)
	if err != nil {
		errorLoggerBroadcast.Printf("Message sending error")
		return
	}
	infoLoggerBroadcast.Printf("Message was sent")
}

func (node *NodeScanner) Run() {
	go handleConnection(node)
	time.Sleep(time.Second * 2)
	for {
		node.BroadcastMSG()
		time.Sleep(time.Second * 30)
	}
}

//func handleBroadcastServer(node *NodeScanner) {
//
//listen, err := net.Listen("tcp", node.AddressR)
//if err != nil {
//	errorLoggerBroadcast.Printf("Server creating error")
//	panic("Server Error")
//}
//defer listen.Close()
//infoLoggerBroadcast.Printf("BroadcastServer was started with %s", node.AddressR)
//for {
//	conn, err := listen.Accept()
//	infoLoggerBroadcast.Printf("New connection")
//	if err != nil {
//		break
//	}
//	go handleConnection(node, conn)
//}
//}

func handleConnection(node *NodeScanner) {
	conn, err := net.ListenPacket("udp4", "9001")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	var (
		buffer  = make([]byte, 512)
		message string
		pack    PackageBroadcast
	)
	for {
		println("Wait")
		length, _, err := conn.ReadFrom(buffer)
		if err != nil {
			break
		}
		message += string(buffer[:length])
		if strings.HasSuffix(message, END_BYTES) {
			message = strings.TrimSuffix(message, END_BYTES)
			break
		}
	}
	err = json.Unmarshal([]byte(message), &pack)
	if err != nil {
		return
	}
	if node.Room == pack.Room && node.login != pack.Login {
		if node.db.GetKey(pack.Login) == "" {
			infoLoggerBroadcast.Printf("Save %s", pack.Login)
			node.db.SetLogin(pack.Login, pack.Key)
			node.db.SetAddress(pack.Key, pack.Address)
			node.SendToAddress(pack.Address)
		} else if node.db.GetKey(pack.Login) != pack.Key {
			infoLoggerBroadcast.Printf("Update data %s", pack.Login)
			node.db.SetLogin(pack.Login, pack.Key)
			node.db.SetAddress(pack.Key, pack.Address)
		} else if node.db.GetAddress(pack.Key) != pack.Address {
			infoLoggerBroadcast.Printf("Update address %s", pack.Login)
			node.db.SetAddress(pack.Key, pack.Address)
		} else {
			errorLoggerBroadcast.Printf("Unknown %s %s", pack.Login, pack.Address)
		}
		infoLoggerBroadcast.Printf("Answer was sent")
		node.SendToAddress(pack.AddressBroadcast)
	} else {
		errorLoggerBroadcast.Printf("Unknown FULL %s %s", pack.Login, pack.Address)
	}

}

func (node *NodeScanner) SendToAddress(address string) {
	var pack = PackageBroadcast{
		Login:            node.login,
		Address:          node.Address,
		AddressBroadcast: node.AddressB,
		Key:              node.Key,
		Room:             node.Room,
	}
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return
	}
	defer conn.Close()
	JSONPack, _ := json.Marshal(pack)
	_, err = conn.Write(JSONPack)
	if err != nil {
		errorLoggerBroadcast.Printf("Message sending error")
		return
	}
	infoLoggerBroadcast.Printf("Message was sent to %s", address)
}
