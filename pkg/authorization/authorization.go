package authorization

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"kubazulo/pkg/utils"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

func AuthorizeRequestDeviceFlow(c utils.AuthorizationConfig) (jsonDf jsonDeviceFlow) {

	formVals := url.Values{}
	formVals.Set("client_id", c.ClientID)
	formVals.Set("scope", utils.CfgClientId+"/.default offline_access")

	uri, _ := url.Parse(utils.AuthorizationURLDevice)
	uri.RawQuery = formVals.Encode()

	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(req)

	if err != nil {
		errors.Wrap(err, "error while trying to get tokens")
	}
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		errors.Wrap(err, "error while trying to read token json body")
	}

	err = json.Unmarshal(body, &jsonDf)
	if err != nil {
		errors.Wrap(err, "error while trying to parse token json body")
	}

	fmt.Fprintln(os.Stderr, "Please visit", jsonDf.VerificationUri, "and put the following code there:", jsonDf.UserCode)
	fmt.Fprintln(os.Stderr, "You have 1 Minute to complete !!")

	if err != nil {
		panic(errors.Wrap(err, "Error while opening DeviceFlow URL"))
	}

	return
}

func AuthorizeRequest(c utils.AuthorizationConfig) (token AuthorizationCode) {
	formVals := url.Values{}
	formVals.Set("grant_type", "authorization_code")
	formVals.Set("redirect_uri", c.RedirectURL())
	formVals.Set("scope", utils.CfgClientId+"/.default")
	formVals.Set("response_type", "code")
	formVals.Set("response_mode", "query")
	formVals.Set("client_id", c.ClientID)
	formVals.Set("state", "12345")
	if strings.ToLower(utils.CfgForceLogin) == "true" {
		formVals.Set("prompt", "login")
	}
	uri, _ := url.Parse(utils.AuthorizationURL)
	uri.RawQuery = formVals.Encode()

	cmd := exec.Command(c.OpenCMD, uri.String())
	err := cmd.Start()
	if err != nil {
		panic(errors.Wrap(err, "Error while opening login URL"))

	}
	running := true
	srv := startLocalListener(c, &token)
	for running {
		if token.Value != "" {
			if err := srv.Shutdown(context.TODO()); err != nil {
				panic(err) // failure/timeout shutting down the server gracefully
			}
			running = false
		}
	}
	return
}
