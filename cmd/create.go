package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Luzifer/expose/ngrok2"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <port>",
	Short: "Create a new HTTPs tunnel to a local port",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Port not specified")
		}

		port, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("Port %q is not a valid integer: %s", args[0], err)
		}

		name := fmt.Sprintf("expose_%d", port)
		if tun, err := client.GetTunnelByName(name); err == nil {
			return fmt.Errorf("Address %q already has an active tunnel with URL %s", tun.Config.Addr, tun.PublicURL)
		}

		inspect, err := cmd.Flags().GetBool("inspect")
		if err != nil {
			return fmt.Errorf("Unable to read inspect flag: %s", err)
		}

		tun, err := client.CreateTunnel(ngrok2.CreateTunnelInput{
			Name:    name,
			Addr:    strconv.Itoa(port),
			Proto:   "http",
			BindTLS: "true",
			Inspect: inspect,
		})

		if err != nil {
			return fmt.Errorf("Unable to create tunnel: %s", err)
		}

		fmt.Printf("Created tunnel to address %q with URL %s\n", tun.Config.Addr, tun.PublicURL)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	createCmd.Flags().BoolP("inspect", "i", false, "Enable HTTP inspection on tunnel")
}
