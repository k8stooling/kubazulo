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
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}

func WriteSession(_LoginMode string, _Expiry int64, _TokenStart int64, _AccessToken string, _RefreshToken string) {
	CreateDirectory(GetHomeDir() + cachepath)
	f, err := os.Create(GetHomeDir() + cachepath + sessionfile)
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
		log.Fatal(err)
	}
	json.Unmarshal([]byte(fileContent), &data)
	return data
}
