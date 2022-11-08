package main

import (
	kn "Diploma/pkg/kernel"
	"fmt"
	"os"
	"time"
)

func main() {
	var node *kn.Client
	var status string
	fmt.Printf("Press r for registate or a for login")
	choice := kn.InputString()
	if choice[0] == 'r' {
		rstatus := kn.Registration()
		if rstatus == 1 {
			println("Successful")
			os.Exit(1)
		} else {
			println("Registration error")
		}
	}
	if choice[0] == 'a' {
		status, node = kn.Authentication()
		println("AUTHENTICATION Status:", status)
		if status != "ok" {
			return
		}
		time.Sleep(time.Second * 15)
	}
	node.ListF2F()
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
