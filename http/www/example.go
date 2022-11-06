package www

import (
	"fmt"
	"github.com/aaronland/go-http-server-tsnet"
	"net/http"
)

func ExampleHandler() http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		who, err := tsnet.GetWhois(req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
		}

		login_name := who.UserProfile.LoginName
		computed_name := who.Node.ComputedName

		msg := fmt.Sprintf("Hello, %s (%s)", login_name, computed_name)
		rsp.Write([]byte(msg))
	}

	h := http.HandlerFunc(fn)
	return h
}
