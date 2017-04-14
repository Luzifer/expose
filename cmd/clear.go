package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove all active tunnels",
	RunE: func(cmd *cobra.Command, args []string) error {
		tunnels, err := client.ListTunnels()
		if err != nil {
			return err
		}

		if len(tunnels) == 0 {
			fmt.Println("No tunnels active to be closed.")
			return nil
		}

		for _, tun := range tunnels {
			if err := client.StopTunnel(tun.Name); err != nil {
				return fmt.Errorf("Unable to stop tunnel %q: %s", tun.Name, err)
			}
			fmt.Printf("Successfully closed tunnel %q to address %q.\n", tun.Name, tun.Config.Addr)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(clearCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clearCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clearCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
