package cmd

import (
	"github.com/spf13/pflag"
)

func Options() FlagOptions {
	return FlagOptions{}
}

type FlagOptions struct {
	ClientID     string
	TenantID     string
	forceLogin   string
	LookBackPort string
	Intermediate string
	ApiEndpoint  string
	LoginMode    string
}

func (o *FlagOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.ClientID, "client-id", "c", "", "Azure Application-ID (Mandatory)")
	fs.StringVarP(&o.TenantID, "tenant-id", "t", "", "Azure Tenant-ID (Mandatory)")
	fs.StringVarP(&o.forceLogin, "force-login", "f", "false", "Re-Usage of Browser Session data")
	fs.StringVarP(&o.LookBackPort, "loopbackport", "l", "58433", "Customize local callback listener")
	fs.StringVarP(&o.Intermediate, "intermediate", "i", "false", "Activate another Token fetcher Endpoint")
	fs.StringVarP(&o.ApiEndpoint, "api-token-endpoint", "a", "", "External Token Endpoint")
	fs.StringVarP(&o.LoginMode, "loginmode", "m", "interactive", "Login Method to be used")
}

func RequiredFlags() []string {
	return []string{"client-id", "tenant-id", "loginmode"}
}

func DependendFlags() []string {
	return []string{"intermediate", "api-token-endpoint"}
}
