package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	kubazulo "main/pkg"
	"os"
	"time"
)

type Spec struct {
	Interactive bool `json:"interactive"`
}

type Status struct {
	ExpirationTimestamp string `json:"expirationTimestamp"`
	Token               string `json:"token"`
}

type Product struct {
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	Spec       Spec   `json:"spec"`
	Status     Status `json:"status"`
}

type CliFlag struct {
	Required bool
	Usage    string
	Name     string
	Address  *string
}

func kubeoutput(accesstoken string) {
	kcoutput := Product{
		Kind:       "ExecCredential",
		ApiVersion: "client.authentication.k8s.io/v1beta1",
		Spec:       Spec{false},
		Status:     Status{kubazulo.ConvertUnixToRFC3339(kubazulo.GetCurrentUnixTime()), accesstoken},
	}
	bytes, _ := json.Marshal(kcoutput)
	fmt.Println(string(bytes))
}

func createNewToken() {
	authConfig := kubazulo.DefaultConfig
	// client id from your setup
	authConfig.ClientID = kubazulo.Cfg_client_id
	// client secret from your setup
	//authConfig.ClientSecret = x.ClientSecret

	// Preform one time login
	authCode := kubazulo.LoginRequest(authConfig)
	t, err := kubazulo.GetTokens(authConfig, authCode, "profile openid offline_access")
	if err != nil {
		panic(err)
	}
	kubeoutput(t.AccessToken)
	kubazulo.WriteSession(kubazulo.GetExpiryUnixTime(int64(t.Expiry)), kubazulo.GetCurrentUnixTime(), t.AccessToken, t.RefreshToken)
}

func main() {
	var _r kubazulo.Session
	var (
		_client_id   string
		_tenant_id   string
		_force_login string
	)

	flag.StringVar(&_client_id, "client-id", "", "client-id missing")
	flag.StringVar(&_tenant_id, "tenant-id", "", "tenant-id missing")
	flag.StringVar(&_force_login, "force-login", "false", "force-login is missing")

	flag.Parse()

	if _client_id == "" || _tenant_id == "" {
		fmt.Println("ERROR: Command can't be executed! \nMissing Mandatory Parameters: (client-id) and (tenant-id)")
		os.Exit(2)
	}

	kubazulo.Cfg_client_id = _client_id
	kubazulo.Cfg_tenant_id = _tenant_id
	kubazulo.Cfg_force_login = _force_login
	kubazulo.FillVariables()

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(home + "/.kube/cache/kubazulo/azuredata.json"); errors.Is(err, os.ErrNotExist) {
		createNewToken()
	} else {
		r := kubazulo.ReadSession()
		_r = r
		if time.Now().Unix() >= _r.ExpirationTimestamp {
			createNewToken()
		} else if time.Now().Unix() <= _r.ExpirationTimestamp && r.RefreshToken != "" {
			t, err := kubazulo.RenewAccessToken(_r.RefreshToken)
			if err != nil {
				log.Fatal(err)
			}
			kubeoutput(t.AccessToken)
			kubazulo.WriteSession(kubazulo.GetExpiryUnixTime(int64(t.Expiry)), kubazulo.GetCurrentUnixTime(), t.AccessToken, t.RefreshToken)
		} else {
			kubeoutput(r.AccessToken)
		}
	}
}