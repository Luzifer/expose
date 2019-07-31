package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Luzifer/expose/ngrok2"
	http_helper "github.com/Luzifer/go_helpers/http"

	"github.com/spf13/cobra"
)

type corsWrap struct {
	next   http.Handler
	header string
}

func newCorsWrap(next http.Handler, header string) *corsWrap { return &corsWrap{next, header} }

func (c corsWrap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", c.header)
	c.next.ServeHTTP(w, r)
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
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

		var handler http.Handler
		handler = http.FileServer(http.Dir("."))
		handler = http_helper.GzipHandler(handler)

		if c, err := cmd.Flags().GetString("cors"); err == nil && c != "" {
			handler = newCorsWrap(handler, c)
		}

		if l, err := cmd.Flags().GetBool("http-logs"); err == nil && l {
			handler = http_helper.NewHTTPLogHandler(handler)
		}

		go http.Serve(listener, handler)
		fmt.Printf("Created HTTP server for this directory with URL %s\nPress Ctrl+C to stop server and tunnel", tun.PublicURL)

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		if d, err := cmd.Flags().GetDuration("timeout"); err == nil && d > 0 {
			fmt.Printf(" (automatic close in %s)", d)
			go func() {
				<-time.After(d)
				c <- os.Interrupt
			}()
		}

		fmt.Println()

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
	RootCmd.AddCommand(serveCmd)

	serveCmd.Flags().DurationP("timeout", "t", 0, "Automatically close tunnel after timeout")
	serveCmd.Flags().String("cors", "", "Enable setting CORS header")
	serveCmd.Flags().Bool("http-logs", true, "Enable HTTP access logs")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
