package gopeer

import (
	"fmt"
	"log"
	"os"
	"time"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func newLogger(object, typeLogger string) *log.Logger {
	logsDir := "./logs/" + object + "/"
	if exists(logsDir) == false {
		os.Mkdir(logsDir, os.ModePerm)
	}
	year, month, day := time.Now().Date()
	filename := fmt.Sprintf("%v-%v-%v.log", day, month.String(), year)
	path := logsDir + filename

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		println("ERROR: Create|Open File", path)
		os.Exit(-1)
	}
	if typeLogger == "ERROR" {
		return log.New(f, typeLogger+": ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		return log.New(f, typeLogger+": ", log.Ldate|log.Ltime)
	}
}
