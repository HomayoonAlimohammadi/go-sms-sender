package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Short: "sms-sender web app",
	Long: `This is SMS Sender web app written in golang.
	Use "serve" command to run the app.`,
	Run: nil,
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringP("config-file", "c", "./config.yaml", "Path to the config file (e.g. ./config.yaml) [Optional]")

	rootCmd.AddCommand(serveCmd, versionCmd)
}
