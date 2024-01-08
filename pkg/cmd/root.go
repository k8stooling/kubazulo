package cmd

import (
	"github.com/spf13/cobra"
	"kubazulo/pkg/token"
	"kubazulo/pkg/utils"
)

func GetToken() *cobra.Command {
	o := Options()
	cmd := &cobra.Command{
		Use:   "get-token",
		Short: "Tries to fetch the Token from AzureAD",
		Long:  `Tries to fetch the Token from AzureAD`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			token.InvokeTokenProcess(cmd.Flags())
		},
	}

	o.AddFlags(cmd.Flags())
	for _, rf := range RequiredFlags() {
		cmd.MarkFlagRequired(rf)
	}

	cmd.Flags().SortFlags = false
	cmd.MarkFlagsRequiredTogether(DependendFlags()...)

	return cmd
}

func Version() *cobra.Command {
	o := Options()
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Gets the current version of kubazulo",
		Long:  `Gets the current version of kubazulo`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			utils.PrintAppVersion()
		},
	}
	o.AddFlags(cmd.Flags())
	return cmd
}
