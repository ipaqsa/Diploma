package kernel

import (
	"encoding/json"
	"net"
	"strconv"
	"strings"
	"time"
)

var infoLoggerBroadcast = newLogger("broadcast", "INFO")
var errorLoggerBroadcast = newLogger("broadcast", "ERROR")

func (client *Client) NewNodeBroadcast(address, login, key string, room uint) *NodeScanner {
	client.DBFriendsInit()
	infoLoggerBroadcast.Printf("Node was created")
	return &NodeScanner{
		Port:        address,
		login:       login,
		db:          client.dbFriends,
		Key:         key,
		Room:        room,
		Connections: make(map[string]string),
	}
}

func (node *NodeScanner) BroadcastMSG() {
	conn, err := net.ListenPacket("udp4", IncrementPortFromAddress(node.Port, 1))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	addr, err := net.ResolveUDPAddr("udp4", "192.168.0.255"+IncrementPortFromAddress(node.Port, 2))
	if err != nil {
		errorLoggerBroadcast.Printf(err.Error())
	}
	var pack = PackageBroadcast{
		Login: node.login,
		Key:   node.Key,
		Room:  node.Room,
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
		time.Sleep(time.Second * 5)
	}
}

func DecrementPortFromAddress(address string) string {
	splited := strings.Split(address, ":")
	port_i, _ := strconv.Atoi(splited[1])
	port_i--
	splited[1] = strconv.Itoa(port_i)
	address_n := strings.Join(splited, ":")
	return address_n
}
func IncrementPortFromAddress(address string, n int) string {
	splited := strings.Split(address, ":")
	port_i, _ := strconv.Atoi(splited[1])
	port_i += n
	splited[1] = strconv.Itoa(port_i)
	address_n := strings.Join(splited, ":")
	return address_n
}

func (node *NodeScanner) PackageAnalysis(message string, pack *PackageBroadcast, addr net.Addr) int {
	err := json.Unmarshal([]byte(message), &pack)
	if err != nil {
		return 0
	}
	if node.Room == pack.Room && node.login != pack.Login {
		address := DecrementPortFromAddress(addr.String())
		if node.db.GetKey(pack.Login) == "" {
			infoLoggerBroadcast.Printf("Save %s", pack.Login)
			node.db.SetLogin(pack.Login, pack.Key)
			node.db.SetAddress(pack.Key, address)
			return 1
		} else if node.db.GetKey(pack.Login) != pack.Key {
			infoLoggerBroadcast.Printf("Update data %s", pack.Login)
			node.db.SetLogin(pack.Login, pack.Key)
			node.db.SetAddress(pack.Key, address)
			return 1
		} else if node.db.GetAddress(pack.Key) != address {
			infoLoggerBroadcast.Printf("Update address %s", pack.Login)
			node.db.SetAddress(pack.Key, address)
			return 1
		} else {
			errorLoggerBroadcast.Printf("Unknown %s %s", pack.Login, addr.String())
		}
	} else {
		errorLoggerBroadcast.Printf("Unknown FULL %s %s", pack.Login, addr.String())
		return 2
	}
	return 1
}

func handleConnection(node *NodeScanner) {
	conn, err := net.ListenPacket("udp4", IncrementPortFromAddress(node.Port, 2))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	var (
		buffer  = make([]byte, 1024)
		message string
		pack    PackageBroadcast
	)
	for {
		length, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			break
		}
		message += string(buffer[:length])
		status := node.PackageAnalysis(message, &pack, addr)
		if status != 0 {
			message = ""
		}
	}

}

func (node *NodeScanner) SendToAddress(address string) {
	var pack = PackageBroadcast{
		Login: node.login,
		Key:   node.Key,
		Room:  node.Room,
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
