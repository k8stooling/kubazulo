package utils

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/spf13/pflag"
)

func GetCurrentUnixTime() int64 {
	return time.Now().Unix()
}

func GetExpiryUnixTime(ExpiryWindow int64) int64 {
	return time.Now().Unix() + ExpiryWindow
}

func ConvertUnixToRFC3339(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format("2006-01-02T15:04:05Z07:00")
}

func GetHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return home
}

func (c AuthorizationConfig) RedirectURL() string {
	if c.RedirectPort != "" {
		c.Host = fmt.Sprintf("%s:%s", c.Host, c.RedirectPort)
	}

	uri := url.URL{
		Host:   c.Host,
		Scheme: c.Scheme,
		Path:   c.RedirectPath,
	}

	return uri.String()
}

func CheckFlagExistence(flags *pflag.FlagSet, name string) bool {
	result, _ := flags.GetString(name)
	if result != "" {
		return true
	} else {
		return false
	}
}

func PrintAppVersion() {
	fmt.Println("Kubazulo Version:", version)
}
