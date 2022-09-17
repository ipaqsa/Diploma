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
	NODE_ADDRESS  = "127.0.0.1:9000"
	BNODE_ADDRESS = "127.0.0.1:9001"
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

func newNode(address, addressBroadcast string, login string) *gp.Client {
	addresses := createAddresses()
	user := gp.LoadUser(login)
	node := gp.NewClient(address, user)
	nodeBroadcast := node.NewNodeBroadcast(address, addressBroadcast, user.Login, node.StringPublic(), user.Room)
	go nodeBroadcast.Run(addresses)
	go gp.NewListener(node).Run(handleFunc)
	return node
}

func AppendFriends(node *gp.Client) {
	for {
		time.Sleep(time.Second * 5)
		node.AppendFriends()
	}
}

func createAddresses() []string {
	var result []string
	pattern := "192.168.0."
	port := ":9001"
	begin := 1
	end := 240
	for ; begin < end; begin++ {
		address := pattern + strconv.Itoa(begin) + port
		result = append(result, address)
	}
	return result
}

func main() {
	user := gp.NewUser(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)), "macbook", "Stefan", "12", 1)
	gp.SaveUser(user)
	node := newNode(NODE_ADDRESS, BNODE_ADDRESS, "macbook")
	go AppendFriends(node)
	time.Sleep(time.Second * 10)
	println(node.GetAllMembers())
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
