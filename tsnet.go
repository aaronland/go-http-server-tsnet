package tsnet

// https://tailscale.com/blog/tsnet-virtual-private-services/
// https://pkg.go.dev/tailscale.com/tsnet

import (
	"context"
	"fmt"
	"github.com/aaronland/go-http-server"
	"net/http"
	"net/url"
	"tailscale.com/tsnet"
)

func init() {
	ctx := context.Background()
	server.RegisterServer(ctx, "tsnet", NewTSNetServer)
}

type TSNetServer struct {
	server.Server
	tsnet_server *tsnet.Server
	hostname     string
	port         string
}

func NewTSNetServer(ctx context.Context, uri string) (server.Server, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}
	
	hostname := u.Hostname()
	port := u.Port()

	q := u.Query()
	
	tsnet_server := &tsnet.Server{
		Hostname: u.Hostname(),
	}

	auth_key := q.Get("auth-key")
	
	if auth_key != "" {
		tsnet_server.AuthKey = auth_key
	}
	
	s := &TSNetServer{
		tsnet_server: tsnet_server,
		hostname:     hostname,
		port:         port,
	}

	return s, nil
}

func (s *TSNetServer) ListenAndServe(ctx context.Context, mux http.Handler) error {

	ln, err := s.tsnet_server.Listen("tcp", s.Address())

	if err != nil {
		return fmt.Errorf("Failed to announce server, %w", err)
	}

	defer ln.Close()

	err = http.Serve(ln, mux)

	if err != nil {
		return fmt.Errorf("Failed to serve requests, %w", err)
	}

	return nil
}

func (s *TSNetServer) Address() string {
	return fmt.Sprintf("%s:%s", s.hostname, s.port)
}
