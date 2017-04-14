package cmd

import (
	"fmt"
	"os"

	"github.com/Luzifer/expose/ngrok2"
	"github.com/spf13/cobra"
)

var client *ngrok2.Client

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "expose",
	Short: "Control ngrok dameon with simple CLI commands",

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		client = ngrok2.New()

		baseURL, err := cmd.Flags().GetString("api-base")
		if err != nil {
			return err
		}
		client.BaseURL = baseURL

		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.expose.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.PersistentFlags().String("api-base", "http://localhost:4040", "Base URL to contact the ngrok daemon")
}
