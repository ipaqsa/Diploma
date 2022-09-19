package kernel

import (
	"bytes"
	"crypto/rsa"
	"errors"
	"math/big"
	"net"
	"sync"
	"time"
)

var infoLogger = newLogger("client", "INFO")
var errorLogger = newLogger("client", "ERROR")

func NewUser(priv *rsa.PrivateKey, login, name string, password string, room uint) *User {
	//if uniqueLogin(login) != true {
	//	return nil
	//}
	pswd := HashSum([]byte(password))
	infoLogger.Printf("New user: %s", login)
	return &User{
		Name:       name,
		Login:      login,
		Password:   pswd,
		Room:       room,
		PrivateKey: priv,
	}
}
func NewClient(address string, user *User) *Client {
	infoLogger.Printf("New node with %s", address)
	return &Client{
		user:        user,
		address:     address,
		mutex:       new(sync.Mutex),
		mapping:     make(map[string]bool),
		connections: make(map[net.Conn]string),
		actions:     make(map[string]chan string),
		f2f:         make(map[string]*rsa.PublicKey),
		f2f_d:       make(map[*rsa.PublicKey]string),
	}
}

func (client *Client) SendMessageTo(login string, pack *Package) (string, error) {
	s := client.InF2F(login)
	if s == false {
		errorLogger.Printf("Client is not found in F2F: %s", login)
		return "", nil
	}
	var (
		err    error
		result string
		hash   = HashPublic(client.f2f[login])
	)
	client.actions[hash] = make(chan string)
	defer delete(client.actions, hash)
	client.send(client.f2f[login], pack)
	select {
	case result = <-client.actions[hash]:
	case <-time.After(time.Duration(settings.WAIT_TIME) * time.Second):
		err = errors.New("time is over")
	}
	if err == nil {
		AddMessage(client.user.Login, login, pack)
	}
	return result, err
}

func (client *Client) Connect(login string, handle func(*Client, *Package)) error {
	key := client.dbFriends.GetKey(login)
	if key == "" {
		errorLogger.Printf("Key is not found %s", login)
		return nil
	}
	address := client.dbFriends.GetAddress(key)
	if address == "" {
		errorLogger.Printf("Address is not found")
		return nil
	}
	conn, err := net.Dial("tcp", address)
	if err != nil {
		errorLogger.Printf("Connection error")
		return err
	}
	infoLogger.Printf("Successful connect to %s", login)
	client.connections[conn] = address
	go handleConn(conn, client, handle)
	return nil
}

//func (client *Client) BroadCastLocal() {
//	listenAddr, err := net.ResolveUDPAddr("udp4", ":8827")
//	if err != nil {
//		panic(err)
//	}
//	list, err := net.ListenUDP("udp4", listenAddr)
//	if err != nil {
//		panic(err)
//	}
//	defer list.Close()
//
//	addr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:9000")
//	if err != nil {
//		panic(err)
//	}
//	_, err = list.WriteTo([]byte("data to transmit"), addr)
//	infoLogger.Printf("Broadcast message was sent")
//	if err != nil {
//		panic(err)
//	}
//}

func (client *Client) Disconnect(address string) {
	for conn, addr := range client.connections {
		if addr == address {
			delete(client.connections, conn)
			infoLogger.Printf("Successful disconnect")
			conn.Close()
		}
	}
}

func (client *Client) Public() *rsa.PublicKey {
	return &client.user.PrivateKey.PublicKey
}
func (client *Client) Private() *rsa.PrivateKey {
	return client.user.PrivateKey
}
func (client *Client) StringPublic() string {
	return StringPublic(&client.user.PrivateKey.PublicKey)
}
func (client *Client) StringPrivate() string {
	return StringPrivate(client.user.PrivateKey)
}
func (client *Client) HashPublic() string {
	return HashPublic(&client.user.PrivateKey.PublicKey)
}

func (client *Client) InF2F(login string) bool {
	if _, ok := client.f2f[login]; ok {
		return true
	}
	return false
}

func (client *Client) ListF2F() []*rsa.PublicKey {
	var list []*rsa.PublicKey
	for _, pub := range client.f2f {
		list = append(list, pub)
	}
	return list
}

func (client *Client) ListF2FAddress() []string {
	var list []string
	for _, address := range client.f2f_d {
		list = append(list, address)
	}
	return list
}

func (client *Client) AppendFriends() {
	members := client.GetAllMembers()
	if members == nil {
		return
	}
	for _, login := range members {
		if _, ok := client.f2f[login]; ok {
			continue
		}
		key := client.dbFriends.GetKey(login)
		if key == "" {
			errorLogger.Printf("Key is not found", login)
			return
		}

		address := client.dbFriends.GetAddress(key)
		if address == "" {
			errorLogger.Printf("Address is not found")
			return
		}
		CreateDialog(client.user.Login, login)
		infoLogger.Printf("Dialog is was created")
		client.f2f[login] = ParsePublic(key)
		client.f2f_d[ParsePublic(key)] = address
		infoLogger.Printf("%s add to %s F2F", login, address)
	}
}

func (client *Client) RemoveFriend(login string) {
	delete(client.f2f, login)
}

func (client *Client) RemoveFriendAddress(pub *rsa.PublicKey) {
	delete(client.f2f_d, pub)
}

func (client *Client) UpdateAddress(pub *rsa.PublicKey, address string) {
	client.f2f_d[pub] = address
}

func (client *Client) send(receiver *rsa.PublicKey, pack *Package) {
	encPack := client.encrypt(receiver, pack)
	bytesPack := EncodePackage(encPack)
	client.mapping[encPack.Body.Hash] = true
	for cn := range client.connections {
		cn.Write(bytes.Join(
			[][]byte{
				[]byte(bytesPack),
				[]byte(settings.END_BYTES),
			},
			[]byte{},
		))
	}
	infoLogger.Printf("Message was sent")
}

func (client *Client) redirect(pack *Package, sender net.Conn) {
	bytesPack := EncodePackage(pack)
	for cn := range client.connections {
		if cn == sender {
			continue
		}
		cn.Write(bytes.Join(
			[][]byte{
				[]byte(bytesPack),
				[]byte(settings.END_BYTES),
			},
			[]byte{},
		))
	}
}

func (client *Client) response(pub *rsa.PublicKey, data string) {
	hash := HashPublic(pub)
	if _, ok := client.actions[hash]; ok {
		client.actions[hash] <- data
	}
}

func (client *Client) encrypt(receiver *rsa.PublicKey, pack *Package) *Package {
	var (
		session = GenerateBytes(settings.SKEY_SIZE)
		rand    = GenerateBytes(settings.RAND_SIZE)
		hash    = HashSum(bytes.Join(
			[][]byte{
				rand,
				Base64Decode(client.StringPublic()),
				Base64Decode(StringPublic(receiver)),
				[]byte(pack.Head.Title),
				[]byte(pack.Body.Data),
				[]byte(pack.Body.Date),
			},
			[]byte{},
		))
		sign = Sign(client.user.PrivateKey, hash)
	)
	return &Package{
		Head: HeadPackage{
			Rand:    Base64Encode(EncryptAES(session, rand)),
			Title:   Base64Encode(EncryptAES(session, []byte(pack.Head.Title))),
			Sender:  Base64Encode(EncryptAES(session, Base64Decode(client.StringPublic()))),
			Session: Base64Encode(EncryptRSA(receiver, session)),
		},
		Body: BodyPackage{
			Date: Base64Encode(EncryptAES(session, []byte(pack.Body.Date))),
			Data: Base64Encode(EncryptAES(session, []byte(pack.Body.Data))),
			Hash: Base64Encode(hash),
			Sign: Base64Encode(sign),
		}}
}

func (client *Client) decrypt(pack *Package) *Package {
	session := DecryptRSA(client.user.PrivateKey, Base64Decode(pack.Head.Session))
	if session == nil {
		return nil
	}
	publicBytes := DecryptAES(session, Base64Decode(pack.Head.Sender))
	if publicBytes == nil {
		return nil
	}
	public := ParsePublic(Base64Encode(publicBytes))
	if public == nil {
		return nil
	}
	size := big.NewInt(1)
	size.Lsh(size, uint(settings.AKEY_SIZE-1))
	if public.N.Cmp(size) == -1 {
		return nil
	}
	titleBytes := DecryptAES(session, Base64Decode(pack.Head.Title))
	if titleBytes == nil {
		return nil
	}
	dataBytes := DecryptAES(session, Base64Decode(pack.Body.Data))
	if dataBytes == nil {
		return nil
	}
	dateBytes := DecryptAES(session, Base64Decode(pack.Body.Date))
	if dateBytes == nil {
		return nil
	}
	rand := DecryptAES(session, Base64Decode(pack.Head.Rand))
	hash := HashSum(bytes.Join(
		[][]byte{
			rand,
			publicBytes,
			Base64Decode(client.StringPublic()),
			titleBytes,
			dataBytes,
			dateBytes,
		},
		[]byte{},
	))
	if Base64Encode(hash) != pack.Body.Hash {
		return nil
	}
	err := Verify(public, hash, Base64Decode(pack.Body.Sign))
	if err != nil {
		return nil
	}
	return &Package{
		Head: HeadPackage{
			Rand:    Base64Encode(rand),
			Title:   string(titleBytes),
			Sender:  Base64Encode(publicBytes),
			Session: Base64Encode(session)},
		Body: BodyPackage{
			Data: string(dataBytes),
			Hash: pack.Body.Hash,
			Sign: pack.Body.Sign,
		},
	}
}

func (client *Client) GetUserINFO() *User {
	return client.user
}

//BroadCast IsUser, UpdateUser, UpdateDialog
