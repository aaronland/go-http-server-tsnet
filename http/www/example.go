package www

import (
	"fmt"
	"net/http"

	"github.com/aaronland/go-http-server-tsnet"
)

// ExampleHandler provides an `http.Handler` that print the login and computed name of
// the Tailscale user invoking the handler.
func ExampleHandler() http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		who, err := tsnet.GetWhoIs(req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		if who.UserProfile == nil {
			http.Error(rsp, "Forbidden", http.StatusForbidden)
			return
		}

		login_name := who.UserProfile.LoginName
		computed_name := who.Node.ComputedName

		msg := fmt.Sprintf("Hello, %s (%s)", login_name, computed_name)
		rsp.Write([]byte(msg))
	}

	h := http.HandlerFunc(fn)
	return h
}
