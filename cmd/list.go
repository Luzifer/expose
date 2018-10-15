package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all currently active tunnels",
	RunE: func(cmd *cobra.Command, args []string) error {
		tunnels, err := client.ListTunnels()
		if err != nil {
			return err
		}

		if j, _ := cmd.Flags().GetBool("json"); j {
			return json.NewEncoder(os.Stdout).Encode(tunnels)
		}

		if len(tunnels) == 0 {
			log.Println("No tunnels are active right now.")
			return nil
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Type", "Address", "Public URL"})
		table.SetBorder(false)

		for _, tun := range tunnels {
			table.Append([]string{
				tun.Name,
				tun.Proto,
				tun.Config.Addr,
				tun.PublicURL,
			})
		}

		table.Render()

		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	listCmd.Flags().Bool("json", false, "Give JSON output instead of table")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
