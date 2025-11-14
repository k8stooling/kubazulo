package utils

import (
	"encoding/json"
	"log"
	"os"

	"github.com/pkg/errors"
)

const (
	sessionfile = "azuredata.json"
	cachepath   = "/.kube/cache/kubazulo/"
)

func CreateDirectory(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, 0700)
		if err != nil {
			log.Println(err, " - Please verify if the folder name already exists or is being used by a file!")
		}
	}
}

func WriteSession(_LoginMode string, _Expiry int64, _TokenStart int64, _AccessToken string, _RefreshToken string) {
	CreateDirectory(GetHomeDir() + cachepath)

	//f, err := os.Create(GetHomeDir() + cachepath + sessionfile)
	f, err := os.OpenFile(GetHomeDir()+cachepath+sessionfile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)

	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := Session{
		CreationOrigin:      _LoginMode,
		TokenStartTimestamp: _TokenStart,
		ExpirationTimestamp: _Expiry,
		AccessToken:         _AccessToken,
		RefreshToken:        _RefreshToken,
	}

	aJson, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	f.Write(aJson)
}

func ReadSession() Session {
	data := Session{}
	fileContent, err := os.ReadFile(GetHomeDir() + cachepath + sessionfile)
	if err != nil {
		log.Fatal(err, " - Please check if folder with name \"kubazulo\" name has been created properly and the name is not already in use!")
	}
	json.Unmarshal([]byte(fileContent), &data)
	return data
}
