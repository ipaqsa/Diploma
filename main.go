package main

import (
	gp "NetworkHiddebLake/gopeer"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	TITLE_MESSAGE   = "MESSAGE!"
	A_NODE_ADDRESS  = "127.0.0.1:9000"
	A_BNODE_ADDRESS = "127.0.0.1:9001"
	B_NODE_ADDRESS  = "127.0.0.1:9002"
	B_BNODE_ADDRESS = "127.0.0.1:9003"
	C_NODE_ADDRESS  = "127.0.0.1:9004"
	C_BNODE_ADDRESS = "127.0.0.1:9005"
	D_NODE_ADDRESS  = "127.0.0.1:9006"
	D_BNODE_ADDRESS = "127.0.0.1:9007"
	E_NODE_ADDRESS  = "127.0.0.1:9008"
	E_BNODE_ADDRESS = "127.0.0.1:9009"
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
	addresses := []string{A_BNODE_ADDRESS, B_BNODE_ADDRESS}
	user := gp.LoadUser(login)
	node := gp.NewClient(address, user)
	nodeBroadcast := node.NewNodeBroadcast(address, addressBroadcast, user.Login, node.StringPublic(), user.Room)
	go nodeBroadcast.Run(addresses)
	go gp.NewListener(node).Run(handleFunc)
	return node
}

func main() {
	alice := gp.NewUser(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)), "alice", "Alice", "12", 1)
	bob := gp.NewUser(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)), "bob", "Bob", "12", 1)
	carl := gp.NewUser(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)), "carl", "Carl", "12", 2)
	dominic := gp.NewUser(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)), "dominic", "Dominic", "12", 1)
	egor := gp.NewUser(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)), "egor", "Egor", "12", 2)

	gp.SaveUser(alice)
	gp.SaveUser(bob)
	gp.SaveUser(carl)
	gp.SaveUser(dominic)
	gp.SaveUser(egor)

	aNode := newNode(A_NODE_ADDRESS, A_BNODE_ADDRESS, "alice")
	bNode := newNode(B_NODE_ADDRESS, B_BNODE_ADDRESS, "bob")
	//cNode := newNode(C_NODE_ADDRESS, C_BNODE_ADDRESS, "carl")
	//dNode := newNode(D_NODE_ADDRESS, D_BNODE_ADDRESS, "dominic")
	//eNode := newNode(E_NODE_ADDRESS, E_BNODE_ADDRESS, "egor")

	time.Sleep(time.Second * 5)
	aNode.AppendFriends()
	err := aNode.Connect("bob", handleFunc)
	if err != nil {
		println(err)
		return
	}
	//for {
	pack := createPackage(TITLE_MESSAGE, "Hello")
	res, err := aNode.SendMessageTo("bob", pack)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	println(res)
	println(bNode)
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
