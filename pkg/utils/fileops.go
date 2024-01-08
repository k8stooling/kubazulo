package utils

import (
	"encoding/json"
	"log"
	"os"

	"github.com/pkg/errors"
)

func CreateDirectory(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}

func WriteSession(Expiry int64, TokenStart int64, _AccessToken string, _RefreshToken string) {
	path := "/.kube/cache/kubazulo/"
	CreateDirectory(GetHomeDir() + path)
	f, err := os.Create(GetHomeDir() + path + "azuredata.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := Session{
		TokenStartTimestamp: TokenStart,
		ExpirationTimestamp: Expiry,
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
	fileContent, err := os.ReadFile(GetHomeDir() + "/.kube/cache/kubazulo/azuredata.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(fileContent), &data)
	return data
}
