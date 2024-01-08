package main

import (
	"github.com/spf13/cobra"
	"kubazulo/pkg/cmd"
	"os"
)

func main() {
	var rootCmd = &cobra.Command{Use: "kubazulo"}
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(cmd.GetToken())
	rootCmd.AddCommand(cmd.Version())
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
