// Package cmd contains the cobra cli commands
package cmd

import (
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zwinslett/speed-daemon/strava"
	"github.com/zwinslett/speed-daemon/telegram"
)

var (
	client *strava.Client
	bot    *telegram.Client
)

var rootCmd = &cobra.Command{
	Use:   "speedd",
	Short: "Strava command line interface",
}

func Execute() {
	client = strava.NewClient()
	err := client.SetAccessToken(context.Background())
	if err != nil {
		log.Fatalf("Auth failed %v", err)
	}
	bot = telegram.NewClient()
	rootCmd.AddCommand(activityByIDCmd(), lastActivityCmd(), zonesCmd(), statsCmd(), notifyCmd(), daemonCmd())
	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
