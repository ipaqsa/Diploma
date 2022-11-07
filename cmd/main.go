package main

import (
	kn "Diploma/pkg/kernel"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	NODE_ADDRESS = ":7000"
)

func Registration() int {
	fmt.Printf("Enter login")
	login := InputString()
	fmt.Printf("Enter name")
	name := InputString()
	fmt.Printf("Enter password")
	password := InputString()
	user := kn.NewUser(kn.GeneratePrivate(kn.Get("AKEY_SIZE").(uint)), login, name, password, 1)
	node := kn.NewClient(NODE_ADDRESS, user)
	node.InitAllDB()
	nodeBroadcast := node.NewNodeBroadcast(NODE_ADDRESS, user.Login, node.StringPublic(), user.Room)
	go nodeBroadcast.Run()
	go kn.NewListener(node).Run()
	go node.RegisterDataSender()
	println("Please wait 10 seconds to register")
	time.Sleep(time.Second * 10)
	err := node.SaveUser(user)
	if err != nil {
		println(err.Error())
		return 0
	}
	return 1
}

func Authentication() (string, *kn.Client) {
	fmt.Printf("Enter login")
	login := InputString()
	fmt.Printf("Enter password")
	password := InputString()
	user := &kn.User{
		Login:      login,
		Password:   "",
		Room:       0,
		PrivateKey: nil,
	}
	status := kn.GetUserFromDB(user, password)
	if status == 1 {
		node := kn.NewClient(NODE_ADDRESS, user)
		node.InitAllDB()
		nodeBroadcast := node.NewNodeBroadcast(NODE_ADDRESS, user.Login, node.StringPublic(), user.Room)
		go nodeBroadcast.Run()
		go kn.NewListener(node).Run()
		go node.RegisterDataSender()
		return "OK", node
	} else {
		return "ERROR Authentication", nil
	}
}

func main() {
	fmt.Printf("Press r for registate or a for login")
	choice := InputString()
	if choice[0] == 'r' {
		rstatus := Registration()
		if rstatus == 1 {
			println("Successful")
			choice = "a"
		} else {
			println("Registration error")
		}
	}
	if choice[0] == 'a' {
		status, _ := Authentication()
		println("AUTHENTICATION Status:", status)
		if status != "ok" {
			return
		}
		time.Sleep(time.Second * 15)
	}
	//for {
	//pack := kn.CreatePackage(InputString())
	//pack := createFilePackage("./data/img.jpg")
	//r, err := node.SendMessageTo("linux", pack)
	//if err != nil {
	//	println(err.Error())
	//	return
	//}
	//println(r)
	//}
}

func InputString() string {
	print(":> ")
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Replace(msg, "\n", "", -1)
}
