package main

import (
	kn "NetworkHiddebLake/kernel"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	TITLE_MESSAGE = "MESSAGE!"
	//NODE_ADDRESS  = "192.168.0.104:8000"
	NODE_ADDRESS = ":8000"
)

func createPackage(title string, data string) *kn.Package {
	return &kn.Package{
		Head: kn.HeadPackage{
			Title: title,
		},
		Body: kn.BodyPackage{
			Date: time.Now().Format("2006-01-02 15:04:05"),
			Data: data,
		},
	}
}

func Registration(address string, user *kn.User) {
	node := kn.NewClient(address, user)
	println("Client register")
	node.DBUsersInit()
	println("DB create")
	err := node.SaveUser(user)
	if err != nil {
		return
	}
}

func Authentication(address, login, password string) (string, *kn.Client) {
	user := &kn.User{
		Login:      login,
		Password:   nil,
		Room:       0,
		PrivateKey: nil,
	}
	status := kn.GetUserFromDB(user, password)
	if status == 1 {
		node := kn.NewClient(address, user)
		nodeBroadcast := node.NewNodeBroadcast(address, user.Login, node.StringPublic(), user.Room)
		go kn.NewListener(node).Run(handleFunc)
		go nodeBroadcast.Run()
		return "OK", node
	} else {
		return "ERROR Authentication", nil
	}
}

func AppendFriends(node *kn.Client) {
	for {
		time.Sleep(time.Second * 2)
		node.AppendFriends()
	}
}

func main() {
	println("Registration:")
	println("Enter name:")
	name := InputString()
	println("Enter login")
	login := InputString()
	println("Enter password")
	password := InputString()
	println("Enter room")
	room := InputString()
	room_i, _ := strconv.Atoi(room)
	room_ui := uint(room_i)
	user := kn.NewUser(kn.GeneratePrivate(kn.Get("AKEY_SIZE").(uint)), login, name, password, room_ui)
	Registration(NODE_ADDRESS, user)
	println("Successful registration")
	status, node := Authentication(NODE_ADDRESS, login, password)
	println("AUTHENTICATION Status", status)
	time.Sleep(time.Second * 10)
	println(node)
}

func handleFunc(client *kn.Client, pack *kn.Package) {
	kn.Handle(TITLE_MESSAGE, client, pack, handleMessage)
}

func handleMessage(client *kn.Client, pack *kn.Package) string {
	sender := kn.ParsePublic(pack.Head.Sender)
	fmt.Printf("\n[%s] => '%s'\n:> ", client.GetLoginFromF2f(sender), pack.Body.Data)
	return "ok"
}

func InputString() string {
	print(":> ")
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Replace(msg, "\n", "", -1)
}
