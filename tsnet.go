package tsnet

// https://tailscale.com/blog/tsnet-virtual-private-services/
// https://github.com/tailscale/tailscale/blob/v1.32.2/tsnet/example/tshello/tshello.go
// https://pkg.go.dev/tailscale.com/tsnet

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/aaronland/go-http-server"
	"net/http"
	"net/url"
	"os"
	"tailscale.com/client/tailscale"
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

		// I don't really understand this...
		// 2022/11/06 10:17:17 Authkey is set; but state is NoState. Ignoring authkey. Re-run with TSNET_FORCE_LOGIN=1 to force use of authkey.

		err := os.Setenv("TSNET_FORCE_LOGIN", "1")

		if err != nil {
			return nil, fmt.Errorf("Failed to set TSNET_FORCE_LOGIN environment variable, %w", err)
		}

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

	// It is important to include the hostname here
	addr := fmt.Sprintf(":%s", s.port)

	ln, err := s.tsnet_server.Listen("tcp", addr)

	if err != nil {
		return fmt.Errorf("Failed to announce server, %w", err)
	}

	defer ln.Close()

	// https://testing
	// 2022/11/06 11:01:23 http: TLS handshake error from a.b.c.d:61940: 400 Bad Request: invalid domain

	// https://{IP}
	// 2022/11/06 11:01:46 http: TLS handshake error from a.b.c.d:53770: no SNI ServerName

	if s.port == "443" {

		ln = tls.NewListener(ln, &tls.Config{
			GetCertificate: tailscale.GetCertificate,
		})
	}

	lc, err := s.tsnet_server.LocalClient()

	if err != nil {
		return fmt.Errorf("Failed to create local client, %w", err)
	}

	who_wrapper := func(next http.Handler) http.Handler {

		fn := func(rsp http.ResponseWriter, req *http.Request) {

			who, err := lc.WhoIs(req.Context(), req.RemoteAddr)

			if err != nil {
				http.Error(rsp, err.Error(), 500)
				return
			}

			who_ctx := context.WithValue(ctx, "who", who)
			who_req := req.WithContext(who_ctx)

			next.ServeHTTP(rsp, who_req)
		}

		return http.HandlerFunc(fn)
	}

	err = http.Serve(ln, who_wrapper(mux))

	if err != nil {
		return fmt.Errorf("Failed to serve requests, %w", err)
	}

	return nil
}

func (s *TSNetServer) Address() string {
	return fmt.Sprintf("%s:%s", s.hostname, s.port)
}
