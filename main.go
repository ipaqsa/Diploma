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
	TITLE_MESSAGE  = "MESSAGE!"
	A_NODE_ADDRESS = "127.0.0.1:8000"
	B_NODE_ADDRESS = "127.0.0.1:7000"
	C_NODE_ADDRESS = "192.168.0.104:7000"
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

func newNode(address string, login string) *gp.Client {
	user := gp.LoadUser(login)
	node := gp.NewClient(address, user)
	go gp.NewListener(node).Run(handleFunc)
	return node
}

func main() {
	//stefan := gp.NewUser(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)), "stefan", "St", "12", 2)
	//alice := gp.NewUser(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)), "alice", "St", "12", 2)
	//gp.SaveUser(stefan)
	//gp.SaveUser(alice)
	//println(gp.IsUser("stefan", "13"))

	aNode := newNode(A_NODE_ADDRESS, "stefan")
	bNode := newNode(B_NODE_ADDRESS, "alice")
	aNode.AppendFriend(bNode.Public(), "alice", B_NODE_ADDRESS)
	aNode.Connect(B_NODE_ADDRESS, handleFunc)
	bNode.Connect(A_NODE_ADDRESS, handleFunc)
	for {
		pack := createPackage(TITLE_MESSAGE, InputString())
		res, err := aNode.SendMessageTo("alice", pack)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		println(res)
	}
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
