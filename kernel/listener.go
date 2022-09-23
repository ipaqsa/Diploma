package kernel

import (
	"net"
	"strings"
)

var infoLoggerListener = newLogger("listener", "INFO")
var errorLoggerListener = newLogger("listener", "ERROR")

func Handle(title string, client *Client, pack *Package, handle func(*Client, *Package) string) bool {
	splited := strings.Split(pack.Head.Title, ":")
	switch splited[0] {
	case title:
		public := ParsePublic(pack.Head.Sender)
		client.send(public, &Package{
			Head: HeadPackage{
				Title: "_" + title,
			},
			Body: BodyPackage{
				Data: handle(client, pack),
			},
		})
	case "_" + title:
		client.response(ParsePublic(pack.Head.Sender), pack.Body.Data)

	default:
		return false
	}
	return true
}

func NewListener(client *Client) *Listener {
	return &Listener{
		client: client,
	}
}

func (listener *Listener) Run(handle func(*Client, *Package)) error {
	var err error
	listener.listen, err = net.Listen("tcp", listener.client.address)
	if err != nil {
		return err
	}
	listener.serve(handle)
	return nil
}

func (listener *Listener) serve(handle func(*Client, *Package)) {
	defer listener.listen.Close()
	for {
		conn, err := listener.listen.Accept()
		if err != nil {
			break
		}
		listener.client.connections[conn] = "client"
		go handleConn(conn, listener.client, handle)
	}
}

func handleConn(conn net.Conn, client *Client, handle func(*Client, *Package)) {
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
		handle(client, decPack)
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
