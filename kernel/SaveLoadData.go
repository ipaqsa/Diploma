package kernel

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

//func SaveUser(user *User) {
//	filename := user.Login + ".json"
//	userJson, err := json.Marshal(user)
//	if err != nil {
//		panic(err)
//	}
//	ioutil.WriteFile(filename, userJson, os.ModePerm)
//}

//func LoadUser(login string) *User {
//	filename := login + ".json"
//	plan, _ := ioutil.ReadFile(filename)
//	var data User
//	json.Unmarshal(plan, &data)
//	return &data
//}

type Dialog struct {
	Messages []Package `json:"messages"`
}

func CreateDialog(from, to string) *Dialog {
	dialog := Dialog{
		Messages: make([]Package, 1),
	}
	SaveDialog(from, to, &dialog)
	return &dialog
}

func SaveDialog(from, to string, dialog *Dialog) {
	var (
		one string
		two string
	)
	if from > to {
		one = from
		two = to
	} else {
		one = to
		two = from
	}
	summary := one + "to" + two
	//id := HashSum([]byte(summary))
	for _, msg := range dialog.Messages {
		println(msg.Body.Data)
	}
	dialogJSON, err := json.Marshal(dialog)
	if err != nil {
		panic(err)
	}
	filename := summary + ".json"
	ioutil.WriteFile(filename, dialogJSON, os.ModePerm)
}

func LoadDialog(from, to string) *Dialog {
	var (
		one string
		two string
	)
	if from > to {
		one = from
		two = to
	} else {
		one = to
		two = from
	}
	summary := one + "to" + two
	//id := HashSum([]byte(summary))
	filename := summary + ".json"
	if _, err := os.Stat("/path/to/whatever"); errors.Is(err, os.ErrNotExist) {
		CreateDialog(from, to)
	}
	plan, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var data Dialog
	json.Unmarshal(plan, &data)
	if data.Messages == nil {
		data.Messages = make([]Package, 1)
	}
	return &data
}

func AddMessage(from, to string, pack *Package) {
	dialog := LoadDialog(from, to)
	dialog.Messages = append(dialog.Messages, *pack)
	SaveDialog(from, to, dialog)
}
