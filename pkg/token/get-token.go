package token

import (
	"encoding/json"
	"fmt"
	"kubazulo/pkg/authorization"
	"kubazulo/pkg/utils"
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

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

func kubeOutput(accessToken string) {
	kcOutput := Product{
		Kind:       "ExecCredential",
		ApiVersion: "client.authentication.k8s.io/v1",
		Spec:       Spec{false},
		Status:     Status{utils.ConvertUnixToRFC3339(utils.GetCurrentUnixTime()), accessToken},
	}
	bytes, _ := json.Marshal(kcOutput)
	fmt.Println(string(bytes))
}

func createNewTokenDeviceFlow() {
	authConfig := utils.DefaultConfig
	authConfig.ClientID = utils.CfgClientId
	authConfig.LoginMode = utils.CfgLoginMode
	authCode := authorization.AuthorizeRequestDeviceFlow(authConfig)

	if utils.CfgIntermediate == "true" {
		var data = authorization.JsonData{
			Code:        authCode.DeviceCode,
			RedirectURI: "http://localhost:" + utils.CfgLoopbackport,
			GrantType:   "urn:ietf:params:oauth:grant-type:device_code",
		}

		t, err := authorization.GetTokenDataApi(data)
		if err != nil {
			panic(err)
		}
		kubeOutput(t.AccessToken)
		utils.WriteSession(utils.GetExpiryUnixTime(int64(t.Expiry)), utils.GetCurrentUnixTime(), t.AccessToken, t.RefreshToken)
	} else {
		for i := 0; i < 12; i++ {
			time.Sleep(5 * time.Second)

			t, err := authorization.GetTokensDeviceCode(authConfig, authCode, "profile openid offline_access")
			if err != nil {
				panic(err)
			}
			if t.AccessToken != "" {
				kubeOutput(t.AccessToken)
				utils.WriteSession(utils.GetExpiryUnixTime(int64(t.Expiry)), utils.GetCurrentUnixTime(), t.AccessToken, t.RefreshToken)
				break
			}
		}
	}

}

func createNewToken() {
	authConfig := utils.DefaultConfig
	authConfig.ClientID = utils.CfgClientId
	authConfig.RedirectPort = utils.CfgLoopbackport

	authCode := authorization.AuthorizeRequest(authConfig)

	if utils.CfgIntermediate == "true" {
		var data = authorization.JsonData{
			Code:        authCode.Value,
			RedirectURI: "http://localhost:" + utils.CfgLoopbackport,
			GrantType:   "authorization_code",
		}

		t, err := authorization.GetTokenDataApi(data)
		if err != nil {
			panic(err)
		}
		kubeOutput(t.AccessToken)
		utils.WriteSession(utils.GetExpiryUnixTime(int64(t.Expiry)), utils.GetCurrentUnixTime(), t.AccessToken, t.RefreshToken)
	} else {
		t, err := authorization.GetTokenAuthCode(authConfig, authCode, "profile openid offline_access")
		if err != nil {
			panic(err)
		}
		kubeOutput(t.AccessToken)
		utils.WriteSession(utils.GetExpiryUnixTime(int64(t.Expiry)), utils.GetCurrentUnixTime(), t.AccessToken, t.RefreshToken)
	}
}

func InvokeTokenProcess(flags *pflag.FlagSet) {
	var _r utils.Session

	utils.InfoLogger.Println("Application invoked.")

	if utils.CheckFlagExistence(flags, "client-id") {
		utils.CfgClientId = flags.Lookup("client-id").Value.String()
	}

	if utils.CheckFlagExistence(flags, "tenant-id") {
		utils.CfgTenantId = flags.Lookup("tenant-id").Value.String()
	}

	if utils.CheckFlagExistence(flags, "force-login") {
		utils.CfgForceLogin = flags.Lookup("force-login").Value.String()
	}

	if utils.CheckFlagExistence(flags, "loopbackport") {
		utils.CfgLoopbackport = flags.Lookup("loopbackport").Value.String()
	}

	if utils.CheckFlagExistence(flags, "intermediate") {
		utils.CfgIntermediate = flags.Lookup("intermediate").Value.String()
	}

	if utils.CheckFlagExistence(flags, "api-token-endpoint") {
		utils.CfgApitokenendpoint = flags.Lookup("api-token-endpoint").Value.String()
	}

	if utils.CheckFlagExistence(flags, "loginmode") {
		utils.CfgLoginMode = flags.Lookup("loginmode").Value.String()
	}

	utils.FillVariables()

	if _, err := os.Stat(utils.GetHomeDir() + "/.kube/cache/kubazulo/azuredata.json"); errors.Is(err, os.ErrNotExist) {
		utils.InfoLogger.Println("Cache File does not exist. New AccessToken obtained from Azure-API.")
		if utils.CfgLoginMode != "devicecode" {
			createNewToken()
		} else {
			createNewTokenDeviceFlow()
		}
	} else {
		r := utils.ReadSession()
		_r = r
		if _r.AccessToken == "" {
			utils.InfoLogger.Println("Cache File does not contain an Access-Token. New AccessToken obtained from Azure-API.")
			createNewToken()
		} else if time.Now().Unix() >= _r.ExpirationTimestamp {
			utils.InfoLogger.Println("Cache File exist but AccessToken is expired. New AccessToken obtained from Azure-API via RefreshToken.")
			t, err := authorization.RenewAccessToken(_r.RefreshToken)
			if err != nil {
				log.Fatal(err)
			}
			kubeOutput(t.AccessToken)
			utils.WriteSession(utils.GetExpiryUnixTime(int64(t.Expiry)), utils.GetCurrentUnixTime(), t.AccessToken, t.RefreshToken)
			utils.InfoLogger.Println("Cache File updated with the latest information from Azure-API.")
		} else {
			utils.InfoLogger.Println("Cache File exist. AccessToken taken from cache file.")
			kubeOutput(r.AccessToken)
		}
	}
}
