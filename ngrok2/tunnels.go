package ngrok2

// Tunnel represents one open tunnel in the daemon
type Tunnel struct {
	Name      string `json:"name"`
	URI       string `json:"uri"`
	PublicURL string `json:"public_url"`
	Proto     string `json:"proto"`
	Config    struct {
		Addr    string `json:"addr"`
		Inspect bool   `json:"inspect"`
	} `json:"config"`
	Metrics TunnelMetrics `json:"metrics"`
}

// TunnelMetrics contains details about the tunnel
type TunnelMetrics struct {
	HTTP struct {
		Count  int64   `json:"count"`
		Rate1  float64 `json:"rate1"`
		Rate5  float64 `json:"rate5"`
		Rate15 float64 `json:"rate15"`
		P50    float64 `json:"p50"`
		P90    float64 `json:"p90"`
		P95    float64 `json:"p95"`
		P99    float64 `json:"p99"`
	} `json:"http"`
	Conns struct {
		Count  int64   `json:"count"`
		Gauge  int64   `json:"gauge"`
		Rate1  float64 `json:"rate1"`
		Rate5  float64 `json:"rate5"`
		Rate15 float64 `json:"rate15"`
		P50    float64 `json:"p50"`
		P90    float64 `json:"p90"`
		P95    float64 `json:"p95"`
		P99    float64 `json:"p99"`
	} `json:"conns"`
}

// CreateTunnelInput represents all options possible to be set on a tunnel
type CreateTunnelInput struct {
	// [required, all] Name of the tunnel to create
	Name string `json:"name"`

	// [required, all] tunnel protocol name, one of `http`, `tcp`, `tls`
	Proto string `json:"proto"`

	// [required, all] forward traffic to this local port number or network address
	Addr string `json:"addr"`

	// [all] enable http request inspection
	Inspect bool `json:"inspect,omitempty"`

	// [http] HTTP basic authentication credentials to enforce on tunneled requests `username:password`
	Auth string `json:"auth,omitempty"`

	// [http] Rewrite the HTTP Host header to this value, or `preserve` to leave it unchanged
	HostHeader string `json:"host_header,omitempty"`

	// [http] bind an HTTPS or HTTP endpoint or both `true`, `false`, or `both`
	BindTLS string `json:"bind_tls,omitempty"`

	// [http, tls] subdomain name to request. If unspecified, uses the tunnel name
	Subdomain string `json:"subdomain,omitempty"`

	// [http, tls] hostname to request (requires reserved name and DNS CNAME)
	Hostname string `json:"hostname,omitempty"`

	// [tls] PEM TLS certificate at this path to terminate TLS traffic before forwarding locally
	CRT string `json:"crt,omitempty"`

	// [tls] PEM TLS private key at this path to terminate TLS traffic before forwarding locally
	Key string `json:"key,omitempty"`

	// [tls] PEM TLS certificate authority at this path will verify incoming TLS client connection certificates.
	ClientCAS string `json:"client_cas,omitempty"`

	// [tcp] bind the remote TCP port on the given address
	RemoteAddr string `json:"remote_addr,omitempty"`
}

// ListTunnels returns a slice of all known open tunnels in the daemon
func (c *Client) ListTunnels() ([]Tunnel, error) {
	var envelope = struct {
		Tunnels []Tunnel `json:"tunnels"`
	}{}

	return envelope.Tunnels, c.do("GET", "/api/tunnels", nil, &envelope)
}

// GetTunnelByName retrieves details for a specific tunnel by its name
func (c *Client) GetTunnelByName(name string) (Tunnel, error) {
	resp := Tunnel{}
	return resp, c.do("GET", "/api/tunnels/"+name, nil, &resp)
}

// StopTunnel closes a tunnel by its name
func (c *Client) StopTunnel(name string) error {
	return c.do("DELETE", "/api/tunnels/"+name, nil, nil)
}

// CreateTunnel starts a new tunnel for given protocol and address
// See https://ngrok.com/docs#tunnel-definitions for details
func (c *Client) CreateTunnel(tunnelOpts CreateTunnelInput) (Tunnel, error) {
	resp := Tunnel{}

	return resp, c.do("POST", "/api/tunnels", tunnelOpts, &resp)
}
