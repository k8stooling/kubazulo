package authorization

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"kubazulo/pkg/utils"
	"net/http"
	"net/url"
)

// AuthorizationCode is a value provided after initial successful
// authentication/authorization, it is used to get access/refresh tokens
type AuthorizationCode struct {
	Value string
}

// Tokens holds access and refresh tokens
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiry       int    `json:"expires_in"`
}

type jsonDeviceFlow struct {
	UserCode        string `json:"user_code"`
	DeviceCode      string `json:"device_code"`
	VerificationUri string `json:"verification_uri"`
}

func GetTokenAuthCode(c utils.AuthorizationConfig, authCode AuthorizationCode, scope string) (t Tokens, err error) {
	formVals := url.Values{}
	formVals.Set("code", authCode.Value)
	formVals.Set("grant_type", "authorization_code")
	formVals.Set("redirect_uri", c.RedirectURL())
	formVals.Set("scope", scope)
	if c.ClientSecret != "" {
		formVals.Set("client_secret", c.ClientSecret)
	}
	formVals.Set("client_id", c.ClientID)
	response, err := http.PostForm(utils.TokenURL, formVals)

	if err != nil {
		return t, errors.Wrap(err, "error while trying to get tokens")
	}
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return t, errors.Wrap(err, "error while trying to read token json body")
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		return t, errors.Wrap(err, "error while trying to parse token json body")
	}

	return
}

func GetTokensDeviceCode(c utils.AuthorizationConfig, authCode jsonDeviceFlow, scope string) (t Tokens, err error) {
	formVals := url.Values{}
	formVals.Set("device_code", authCode.DeviceCode)
	formVals.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	formVals.Set("scope", scope)
	if c.ClientSecret != "" {
		formVals.Set("client_secret", c.ClientSecret)
	}
	formVals.Set("client_id", c.ClientID)
	response, err := http.PostForm(utils.TokenURL, formVals)

	if err != nil {
		return t, errors.Wrap(err, "error while trying to get tokens")
	}
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return t, errors.Wrap(err, "error while trying to read token json body")
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		return t, errors.Wrap(err, "error while trying to parse token json body")
	}

	return
}

func RenewAccessToken(refreshToken string) (t Tokens, err error) {
	if utils.CfgIntermediate == "true" {
		var data = JsonData{}
		data.GrantType = "refresh_token"
		data.Code = refreshToken

		t, err := GetTokenDataApi(data)
		return t, err
	} else {
		formVals := url.Values{}
		formVals.Set("refresh_token", refreshToken)
		formVals.Set("grant_type", "refresh_token")
		formVals.Set("client_id", utils.CfgClientId)
		response, err := http.PostForm(utils.TokenURL, formVals)

		if err != nil {
			return t, errors.Wrap(err, "error while trying to get tokens")
		}
		body, err := ioutil.ReadAll(response.Body)

		if err != nil {
			return t, errors.Wrap(err, "error while trying to read token json body")
		}

		err = json.Unmarshal(body, &t)
		if err != nil {
			return t, errors.Wrap(err, "error while trying to parse token json body")
		}

		return t, err
	}
}
