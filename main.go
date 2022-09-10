package main

import (
	gp "NetworkHiddebLake/gopeer"
	"fmt"
	"os"
	"time"
)

const (
	TITLE_MESSAGE  = "MESSAGE!"
	A_NODE_ADDRESS = ":8000"
	B_NODE_ADDRESS = ":9090"
)

func main() {
	aClient := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))
	bClient := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))

	aNode := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))
	bNode := gp.NewClient(gp.GeneratePrivate(gp.Get("AKEY_SIZE").(uint)))

	go gp.NewListener(A_NODE_ADDRESS, aNode).Run(handleFunc)
	go gp.NewListener(B_NODE_ADDRESS, bNode).Run(handleFunc)

	time.Sleep(500 * time.Millisecond)

	aClient.Connect(A_NODE_ADDRESS, handleFunc)
	bClient.Connect(B_NODE_ADDRESS, handleFunc)

	aNode.Connect(B_NODE_ADDRESS, handleFunc)

	res, err := aClient.Send(bClient.Public(), &gp.Package{
		Head: gp.HeadPackage{
			Title: TITLE_MESSAGE,
		},
		Body: gp.BodyPackage{
			Data: "Hello world!",
		},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	println(res)
}

func handleFunc(client *gp.Client, pack *gp.Package) {
	gp.Handle(TITLE_MESSAGE, client, pack, handleMessage)
}

func handleMessage(client *gp.Client, pack *gp.Package) string {
	sender := gp.ParsePublic(pack.Head.Sender)
	fmt.Printf("[%s] => '%s'\n", gp.HashPublic(sender), pack.Body.Data)
	return "ok"
}
