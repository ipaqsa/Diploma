package main

import (
	kn "Diploma/kernel"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	TITLE_MESSAGE = "MSG"
	TITLE_FILE    = "FILE:"
	NODE_ADDRESS  = ":7000"
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
func createFilePackage(path string) *kn.Package {
	splited := strings.Split(path, "/")
	filename := splited[len(splited)-1]
	bytes, _ := kn.GetFileBytes(path)
	return &kn.Package{
		Head: kn.HeadPackage{
			Title: TITLE_FILE + filename,
		},
		Body: kn.BodyPackage{
			Date: time.Now().Format("2006-01-02 15:04:05"),
			Data: kn.Base64Encode(bytes),
		},
	}
}

func Registration(address string, user *kn.User) {
	node := kn.NewClient(address, user)
	node.DBUsersInit()
	node.DBDialogsInit()
	err := node.SaveUser(user)
	if err != nil {
		return
	}
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
		go kn.NewListener(node).Run(handleFileFunc)
		go nodeBroadcast.Run()
		go AppendFriends(node)
		return "OK", node
	} else {
		return "ERROR Authentication", nil
	}
}

func AppendFriends(node *kn.Client) {
	for {
		time.Sleep(time.Second * 5)
		node.AppendFriends()
	}
}

func main() {
	user := kn.NewUser(kn.GeneratePrivate(kn.Get("AKEY_SIZE").(uint)), "mac", "Stefan", "12", 1)
	Registration(NODE_ADDRESS, user)
	status, node := Authentication(NODE_ADDRESS, "mac", "12")
	println("AUTHENTICATION Status", status)
	time.Sleep(time.Second * 25)
	//for {
	//pack := createPackage(TITLE_MESSAGE, InputString())
	pack := createFilePackage("./data/img.jpg")
	r, err := node.SendMessageTo("linux", pack, handleFileFunc)
	if err != nil {
		println(err.Error())
		return
	}
	println(r)
	//}
}

func handleMSGFunc(client *kn.Client, pack *kn.Package) {
	kn.Handle(TITLE_MESSAGE, client, pack, handleMessage)
}

func handleFileFunc(client *kn.Client, pack *kn.Package) {
	kn.Handle(TITLE_FILE, client, pack, handleFile)
}

func handleMessage(client *kn.Client, pack *kn.Package) string {
	dialogName := kn.GetDialogName(client.GetLogin(pack.Head.Sender), client.GetUserINFO().Login)
	client.AddMessage(dialogName, pack)
	fmt.Printf("\n[%s] => '%s'\n:> ", client.GetLogin(pack.Head.Sender), pack.Body.Data)
	return "ok"
}
func handleFile(client *kn.Client, pack *kn.Package) string {
	filename := strings.Split(pack.Head.Title, ":")[1]
	kn.SaveFileFromByte("./data/"+filename, kn.Base64Decode(pack.Body.Data))
	return "ok"
}

func InputString() string {
	print(":> ")
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Replace(msg, "\n", "", -1)
}
