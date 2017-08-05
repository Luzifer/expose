package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/net/webdav"

	"github.com/Luzifer/expose/ngrok2"
	http_helper "github.com/Luzifer/go_helpers/http"
	"github.com/spf13/cobra"
)

// webdavCmd represents the serve command
var webdavCmd = &cobra.Command{
	Use:   "webdav",
	Short: "Exposes the current directory using WebDAV",
	RunE: func(cmd *cobra.Command, args []string) error {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return fmt.Errorf("Unable to bind local HTTP server: %s", err)
		}

		_, port, err := net.SplitHostPort(listener.Addr().String())
		if err != nil {
			return fmt.Errorf("Unable to determine bound port: %s", err)
		}

		name := fmt.Sprintf("expose_%s", port)
		if tun, err := client.GetTunnelByName(name); err == nil {
			return fmt.Errorf("Address %q already has an active tunnel with URL %s", tun.Config.Addr, tun.PublicURL)
		}

		tun, err := client.CreateTunnel(ngrok2.CreateTunnelInput{
			Name:    name,
			Addr:    port,
			Proto:   "http",
			BindTLS: "true",
			Inspect: false,
		})

		if err != nil {
			return fmt.Errorf("Unable to create tunnel: %s", err)
		}

		srv := &webdav.Handler{
			Prefix:     "/",
			FileSystem: webdav.Dir("."),
			LockSystem: webdav.NewMemLS(),
		}

		logging := http_helper.NewHTTPLogHandler(srv)

		go http.Serve(listener, logging)
		fmt.Printf("Created HTTP server for this directory with URL %s\nPress Ctrl+C to stop server and tunnel", tun.PublicURL)

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		for range c {
			if err := client.StopTunnel(tun.Name); err != nil {
				return fmt.Errorf("Unable to stop tunnel %q: %s", tun.Name, err)
			}
			fmt.Printf("\nSuccessfully stopped server and tunnel.\n")
			return nil
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(webdavCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// webdavCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// webdavCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
