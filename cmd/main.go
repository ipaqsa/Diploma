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
	node.DBUsersInit()
	node.DBDialogsInit()
	node.DBHashesInit()
	go node.RegisterDataSender()
	println("Please wait 10 seconds to register")
	time.Sleep(time.Second * 10)
	err := node.SaveUser(user)
	if err != nil {
		return 0
	}
	return 1
}

func Authentication(address, login, password string) (string, *kn.Client) {
	user := &kn.User{
		Login:      login,
		Password:   "",
		Room:       0,
		PrivateKey: nil,
	}
	status := kn.GetUserFromDB(user, password)
	if status == 1 {
		node := kn.NewClient(address, user)
		node.DBUsersInit()
		node.DBDialogsInit()
		nodeBroadcast := node.NewNodeBroadcast(address, user.Login, node.StringPublic(), user.Room)
		go node.RegisterDataSender()
		go kn.NewListener(node).Run()
		go nodeBroadcast.Run()
		return "OK", node
	} else {
		return "ERROR Authentication", nil
	}
}

func main() {
	rstatus := Registration()
	if rstatus == 1 {
		println("Successful")
	} else {
		println("Registration error")
	}
	time.Sleep(time.Second * 15)
	//status, node := Authentication(NODE_ADDRESS, "mac", "12")
	//println("AUTHENTICATION Status", status)
	//time.Sleep(time.Second * 15)
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
