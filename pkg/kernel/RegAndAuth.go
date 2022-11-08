package kernel

import (
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

	user := NewUser(GeneratePrivate(Get("AKEY_SIZE").(uint)), login, name, password, 1)
	node := NewClient(NODE_ADDRESS, user)

	node.InitAllDB()
	nodeBroadcast := node.NewNodeBroadcast(NODE_ADDRESS, user.Login, node.StringPublic(), user.Room)

	go nodeBroadcast.Run()
	go NewListener(node).Run()
	go node.RegisterDataSender()
	time.Sleep(time.Second * 3)

	err := node.SaveUser(user)
	if err != nil {
		println(err.Error())
		return 0
	}
	return 1
}

func Authentication() (string, *Client) {
	fmt.Printf("Enter login")
	login := InputString()
	fmt.Printf("Enter password")
	password := InputString()

	user := &User{
		Login:      login,
		Password:   password,
		Room:       0,
		PrivateKey: GeneratePrivate(Get("AKEY_SIZE").(uint)),
	}
	status := GetUserFromDB(user)
	node := NewClient(NODE_ADDRESS, user)
	if status == 0 {
		err := node.SaveUser(user)
		if err != nil {
			println(err.Error())
			return "", nil
		}
	}
	node.InitAllDB()
	nodeBroadcast := node.NewNodeBroadcast(NODE_ADDRESS, user.Login, node.StringPublic(), user.Room)
	go nodeBroadcast.HandleConnection()
	go NewListener(node).Run()

	status2th := node.AuthenticationDataSender()
	println(status2th)
	if status2th == 1 || (status == 1 && status2th == -1) {
		go nodeBroadcast.Run()
		go node.RegisterDataSender()
		return "ok", node
	} else {
		return "error", nil
	}
}

func InputString() string {
	print(":> ")
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Replace(msg, "\n", "", -1)
}
