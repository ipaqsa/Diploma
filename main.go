package main

import (
	kn "NetworkHiddebLake/kernel"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	TITLE_MESSAGE = "MESSAGE!"
	//NODE_ADDRESS  = "192.168.0.104:8000"
	NODE_ADDRESS = ":8000"
	//BNODE_ADDRESS = "192.168.0.104:8001"
	BNODE_ADDRESS = ":8001"
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
	node.DBUsersInit()
	err := node.SaveUser(user)
	if err != nil {
		return
	}
}

func Authentication(address, addressBroadcast, login, password string) (string, *kn.Client) {
	user := &kn.User{
		Login:      login,
		Password:   nil,
		Room:       0,
		PrivateKey: nil,
	}
	status := kn.GetUserFromDB(user, password)
	if status == 1 {
		addresses := createAddresses()
		node := kn.NewClient(address, user)
		nodeBroadcast := node.NewNodeBroadcast(address, addressBroadcast, user.Login, node.StringPublic(), user.Room)
		go kn.NewListener(node).Run(handleFunc)
		go nodeBroadcast.Run(addresses)
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

func createAddresses() []string {
	var result []string
	//pattern := "192.168.0."
	//port := ":8001"
	//begin := 140
	//end := 240
	//for ; begin < end; begin++ {
	//	address := pattern + strconv.Itoa(begin) + port
	//	if address != BNODE_ADDRESS {
	//		result = append(result, address)
	//	}
	//}
	result = append(result, ":9001")
	return result
}

func main() {
	user := kn.NewUser(kn.GeneratePrivate(kn.Get("AKEY_SIZE").(uint)), "linuxFrom8000", "Stefan", "12", 1)
	Registration(NODE_ADDRESS, user)
	status, node := Authentication(NODE_ADDRESS, BNODE_ADDRESS, "linuxFrom8000", "12")
	println(status)
	if node != nil {
		go AppendFriends(node)
		countf := len(node.ListF2F())
		for countf < 1 {
			countf = len(node.ListF2F())
			time.Sleep(time.Second * 5)
		}
		err := node.Connect(node.ListF2F()[0], handleFunc)
		if err != nil {
			println(err)
			return
		}
		for {
			pack := createPackage(TITLE_MESSAGE, InputString())
			_, err := node.SendMessageTo(node.ListF2F()[0], pack)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
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
