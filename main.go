package main

import (
	gp "NetworkHiddebLake/gopeer"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	TITLE_MESSAGE = "MESSAGE!"
	NODE_ADDRESS  = "192.168.0.104:8000"
	BNODE_ADDRESS = "192.168.0.104:8001"
)

func createPackage(title string, data string) *gp.Package {
	return &gp.Package{
		Head: gp.HeadPackage{
			Title: title,
		},
		Body: gp.BodyPackage{
			Date: time.Now().Format("2006-01-02 15:04:05"),
			Data: data,
		},
	}
}

func Registration(address string, user *gp.User) {
	//addresses := createAddresses()
	//user := gp.LoadUser(login)
	node := gp.NewClient(address, user)
	node.DBUsersInit()
	node.SaveUser(user)
	//nodeBroadcast := node.NewNodeBroadcast(address, addressBroadcast, user.Login, node.StringPublic(), user.Room)
	//go nodeBroadcast.Run(addresses)
	//return node
}

func Authentication(address, addressBroadcast, login, password string) string {
	user := &gp.User{
		Login:      login,
		Password:   nil,
		Room:       0,
		PrivateKey: nil,
	}
	status := gp.GetUserFromDB(user, password)
	if status == 1 {
		addresses := createAddresses()
		node := gp.NewClient(address, user)
		nodeBroadcast := node.NewNodeBroadcast(address, addressBroadcast, user.Login, node.StringPublic(), user.Room)
		go nodeBroadcast.Run(addresses)
	}
	return "ERROR Authentication"
}

func AppendFriends(node *gp.Client) {
	for {
		time.Sleep(time.Second * 2)
		node.AppendFriends()
	}
}

func createAddresses() []string {
	var result []string
	pattern := "192.168.0."
	port := ":8001"
	begin := 140
	end := 240
	for ; begin < end; begin++ {
		address := pattern + strconv.Itoa(begin) + port
		if address != BNODE_ADDRESS {
			result = append(result, address)
		}
	}
	return result
}

func main() {
	//user := gp.NewUser(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)), "linux", "Stefan", "12", 1)
	//Registration(NODE_ADDRESS, user)
	Authentication(NODE_ADDRESS, BNODE_ADDRESS, "linux", "12")
	//go AppendFriends(node)
	//countf := len(node.ListF2F())
	//for countf < 1 {
	//	countf = len(node.ListF2F())
	//	time.Sleep(time.Second * 5)
	//}
	//gp.NewListener(node).Run(handleFunc)
	//err := node.Connect("bob", handleFunc)
	//if err != nil {
	//	println(err)
	//	return
	//}
	//for {
	//pack := createPackage(TITLE_MESSAGE, "Hello")
	//res, err := aNode.SendMessageTo("bob", pack)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//println(res)
	//println(bNode)
	//}
}

func handleFunc(client *gp.Client, pack *gp.Package) {
	gp.Handle(TITLE_MESSAGE, client, pack, handleMessage)
}

func handleMessage(client *gp.Client, pack *gp.Package) string {
	sender := gp.ParsePublic(pack.Head.Sender)
	fmt.Printf("[%s] => '%s'\n", gp.HashPublic(sender), pack.Body.Data)
	return "ok"
}

func InputString() string {
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Replace(msg, "\n", "", -1)
}
