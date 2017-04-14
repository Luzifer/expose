package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <port>",
	Short: "Close tunnel to specified local port",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Port not specified")
		}

		port, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("Port %q is not a valid integer: %s", args[0], err)
		}

		name := fmt.Sprintf("expose_%d", port)
		tun, err := client.GetTunnelByName(name)
		if err != nil {
			return fmt.Errorf("There is no active tunnel for port %d", port)
		}

		if err := client.StopTunnel(tun.Name); err != nil {
			return fmt.Errorf("Unable to stop tunnel %q: %s", tun.Name, err)
		}
		fmt.Printf("Successfully closed tunnel %q to address %q.\n", tun.Name, tun.Config.Addr)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
