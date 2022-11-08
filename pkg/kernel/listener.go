package kernel

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

const (
	TITLE_MESSAGE        = "MSG"
	TITLE_FILE           = "FILE"
	TITLE_REGISTRATION   = "REGISTER"
	TITLE_AUTHENTICATION = "AUTH"
)

var infoLoggerListener = newLogger("listener", "INFO")
var errorLoggerListener = newLogger("listener", "ERROR")

func Handle(client *Client, pack *Package) bool {
	splited := strings.Split(pack.Head.Title, ":")
	switch splited[0] {
	case TITLE_AUTHENTICATION:
		infoLoggerListener.Printf("AuthPack was received")
		public := ParsePublic(pack.Head.Sender)
		client.send(public, &Package{
			Head: HeadPackage{
				Title: "_" + splited[0],
			},
			Body: BodyPackage{
				Data: handleAuthentication(client, pack),
			},
		})
	case TITLE_REGISTRATION:
		infoLoggerListener.Printf("RegisterPack was received")
		public := ParsePublic(pack.Head.Sender)
		client.send(public, &Package{
			Head: HeadPackage{
				Title: "_" + splited[0],
			},
			Body: BodyPackage{
				Data: handleRegistration(client, pack),
			},
		})
	case TITLE_MESSAGE:
		infoLoggerListener.Printf("MsgPack was received")
		public := ParsePublic(pack.Head.Sender)
		client.send(public, &Package{
			Head: HeadPackage{
				Title: "_" + splited[0],
			},
			Body: BodyPackage{
				Data: handleMessage(client, pack),
			},
		})
	case TITLE_FILE:
		infoLoggerListener.Printf("FilePack was received")
		public := ParsePublic(pack.Head.Sender)
		client.send(public, &Package{
			Head: HeadPackage{
				Title: "_" + splited[0],
			},
			Body: BodyPackage{
				Data: handleFile(client, pack),
			},
		})
	case "_" + TITLE_MESSAGE:
	case "_" + TITLE_FILE:
	case "_" + TITLE_REGISTRATION:
	case "_" + TITLE_AUTHENTICATION:
		client.response(ParsePublic(pack.Head.Sender), pack.Body.Data)

	default:
		return false
	}
	return true
}
func NewListener(client *Client) *Listener {
	go client.AppendFriends()
	return &Listener{
		client: client,
	}
}

func (listener *Listener) Run() {
	var err error
	listener.listen, err = net.Listen("tcp", listener.client.address)
	if err != nil {
		errorLoggerListener.Printf("Listen error: %s", err.Error())
		return
	}
	listener.serve()
	return
}
func (listener *Listener) serve() {
	defer listener.listen.Close()
	for {
		conn, err := listener.listen.Accept()
		if err != nil {
			break
		}
		listener.client.connections[conn] = "client"
		go handleConn(conn, listener.client)
	}
}
func handleConn(conn net.Conn, client *Client) {
	defer func() {
		conn.Close()
		delete(client.connections, conn)
	}()
	for {
		pack := readPackage(conn)
		if pack == nil {
			break
		}
		client.mutex.Lock()
		if _, ok := client.mapping[pack.Body.Hash]; ok {
			client.mutex.Unlock()
			continue
		}
		if len(client.mapping) >= int(settings.MAPP_SIZE) {
			client.mapping = make(map[string]bool)
		}
		client.mapping[pack.Body.Hash] = true
		client.mutex.Unlock()

		decPack := client.decrypt(pack)
		if decPack == nil {
			client.redirect(pack, conn)
			continue
		}
		client.mutex.Lock()
		client.mutex.Unlock()
		Handle(client, decPack)
	}
}

func readPackage(conn net.Conn) *Package {
	var (
		message string
		size    = uint(0)
		buffer  = make([]byte, settings.BUFF_SIZE)
	)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			return nil
		}
		size += uint(length)
		if size >= settings.PACK_SIZE {
			return nil
		}
		message += string(buffer[:length])
		if strings.Contains(message, settings.END_BYTES) {
			message = strings.Split(message, settings.END_BYTES)[0]
			break
		}
	}
	return DecodePackage(message)
}

func handleMessage(client *Client, pack *Package) string {
	dialogName := GetDialogName(client.GetLogin(pack.Head.Sender), client.GetUserINFO().Login)
	client.AddMessage(dialogName, pack)
	fmt.Printf("\n[%s] => '%s'\n:> ", client.GetLogin(pack.Head.Sender), pack.Body.Data)
	return "ok"
}
func handleFile(client *Client, pack *Package) string {
	filename := strings.Split(pack.Head.Title, ":")[1]
	err := SaveFileFromByte("./data/"+filename, Base64Decode(pack.Body.Data))
	if err != nil {
		return ""
	}
	return "ok"
}
func handleRegistration(client *Client, pack *Package) string {
	login := strings.Split(pack.Head.Title, ":")[1]
	err := client.AddHash(pack.Body.Data, pack.Body.Date, login)
	if err != nil {
		return ""
	}
	return "ok"
}
func handleAuthentication(client *Client, pack *Package) string {
	login := strings.Split(pack.Head.Title, ":")[1]
	hash := client.GetHash(login)
	if hash == "" {
		return "nk"
	}
	var passLogin string
	err := json.Unmarshal(Base64Decode(pack.Body.Data), &passLogin)
	if err != nil {
		return "no"
	}
	hashNew, err := HashSumAuth(passLogin)
	if err != nil {
		return "no"
	}
	if hash == hashNew {
		return "ok"
	}
	return "no"
}
