package kernel

import (
	"encoding/json"
	"strings"
	"time"
)

func CreatePackage(data string) *Package {
	return &Package{
		Head: HeadPackage{
			Title: TITLE_MESSAGE,
		},
		Body: BodyPackage{
			Date: time.Now().Format("2006-01-02 15:04:05"),
			Data: data,
		},
	}
}

func CreateAuthenticationPackage(user *User) *Package {
	userjson, err := json.Marshal(user.Password + user.Login)
	if err != nil {
		return nil
	}
	return &Package{
		Head: HeadPackage{
			Title: TITLE_AUTHENTICATION + ":" + user.Login,
		},
		Body: BodyPackage{
			Date: time.Now().Format("2006-01-02 15:04:05"),
			Data: Base64Encode(userjson),
		},
	}
}

func CreateRegistrationPackage(user *User) *Package {
	userjson, err := json.Marshal(user.Password + user.Login)
	if err != nil {
		return nil
	}
	return &Package{
		Head: HeadPackage{
			Title: TITLE_REGISTRATION + ":" + user.Login,
		},
		Body: BodyPackage{
			Date: time.Now().Format("2006-01-02 15:04:05"),
			Data: Base64Encode(userjson),
		},
	}
}

func CreateFilePackage(path string) *Package {
	splited := strings.Split(path, "/")
	filename := splited[len(splited)-1]
	bytes, _ := GetFileBytes(path)
	return &Package{
		Head: HeadPackage{
			Title: TITLE_FILE + ":" + filename,
		},
		Body: BodyPackage{
			Date: time.Now().Format("2006-01-02 15:04:05"),
			Data: Base64Encode(bytes),
		},
	}
}
