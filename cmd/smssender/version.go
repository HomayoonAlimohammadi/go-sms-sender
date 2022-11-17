package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "smssender version",
	Long:  "SMS Sender web app version",
	Run:   version,
}

func version(cmd *cobra.Command, args []string) {
	fmt.Println("version 1.0.0")
}
