package utils

import (
	"log"
	"os"
)

const msBaseURL = "https://login.microsoftonline.com/"
const version = "0.0.7-beta"

type AuthorizationConfig struct {
	Host         string
	Scheme       string
	RedirectPort string
	RedirectPath string
	Scope        string
	ClientID     string
	OpenCMD      string
	ClientSecret string
	LoginMode    string
}

var DefaultConfig = AuthorizationConfig{
	Host:         "localhost",
	Scheme:       "http",
	RedirectPath: "/",
	Scope:        "openid profile offline_access user.read",
	OpenCMD:      "open",
}

var (
	CfgClientId         string
	CfgTenantId         string
	CfgForceLogin       string
	CfgLoopbackport     string
	CfgIntermediate     string
	CfgApitokenendpoint string
	CfgLoginMode        string
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	DebugLogger   *log.Logger
)

var (
	AuthorizationURL       string
	AuthorizationURLDevice string
	TokenURL               string
)

type Session struct {
	TokenStartTimestamp int64  `json:"tokenstartTimestamp"`
	ExpirationTimestamp int64  `json:"expirationTimestamp"`
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
}

func FillVariables() {
	AuthorizationURL = msBaseURL + CfgTenantId + "/oauth2/v2.0/authorize"
	AuthorizationURLDevice = msBaseURL + CfgTenantId + "/oauth2/v2.0/devicecode"
	TokenURL = msBaseURL + CfgTenantId + "/oauth2/v2.0/token"
}

func init() {
	var logpath = GetHomeDir() + "/.kube/kubazulo/"
	CreateDirectory(logpath)

	file, err := os.OpenFile(logpath+"application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

const SuccessMsg = `
<!DOCTYPE html> <html lang="en"> 
<head>
<style>
body {
  background-image: url('https://wallpaperaccess.com/full/4834955.jpg');
  background-repeat: no-repeat;
  background-attachment: fixed;
  background-size: 100% 100%;
}
</style>
<meta charset="UTF-8">
    <title>Azure Platform Authentication Service</title>
</head>
<body>
	<p style="background-image: url('https://wallpaperaccess.com/full/4834955.jpg');"></p>
    <h2><p style="color: white">You have been successfully authenticated and now ready to communicate with the API-Server</br></br>
	You can close the Browser window now and get back to the command-line!</p></h2>
</body>
</html>
`
