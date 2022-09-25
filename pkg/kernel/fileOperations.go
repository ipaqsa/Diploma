package kernel

import (
	"io"
	"io/ioutil"
	"os"
)

func GetFileBytes(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

func GetHashFromFile(path string) (string, error) {
	data, err := GetFileBytes(path)
	if err != nil {
		return "", err
	}
	return Base64Encode(HashSum(data)), err
}

func EncryptFile(password string, path string) ([]byte, error) {
	psswd, err := paddingPassword(password)
	if err != nil {
		return nil, err
	}
	data, err := GetFileBytes(path)
	if err != nil {
		return nil, err
	}
	return EncryptAES(psswd, data), nil
}

func DecryptFile(password string, encryptData []byte) ([]byte, error) {
	psswd, err := paddingPassword(password)
	if err != nil {
		return nil, err
	}
	data := DecryptAES(psswd, encryptData)
	return data, nil
}

func SaveFileFromByte(path string, data []byte) error {
	err := ioutil.WriteFile(path, data, 0664)
	if err != nil {
		println(err.Error())
		return err
	}
	return nil
}

func EncryptAndSaveFile(password, path string) error {
	encr, err := EncryptFile(password, path)
	if err != nil {
		return err
	}
	err = SaveFileFromByte(path, encr)
	if err != nil {
		return err
	}
	return nil
}

func DecryptAndSaveFile(password, path string) error {
	psswd, err := paddingPassword(password)
	if err != nil {
		return err
	}
	data, err := GetFileBytes(path)
	if err != nil {
		return err
	}
	decr := DecryptAES(psswd, data)
	err = SaveFileFromByte(path, decr)
	if err != nil {
		return err
	}
	return nil
}
