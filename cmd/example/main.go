// example is a command-line application demonstrating how a tsnet-enabled HTTP server works.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/aaronland/go-http-server"
	_ "github.com/aaronland/go-http-server-tsnet"
	"github.com/aaronland/go-http-server-tsnet/http/www"
)

func main() {

	server_uri := flag.String("server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")

	flag.Parse()

	ctx := context.Background()

	s, err := server.NewServer(ctx, *server_uri)

	if err != nil {
		log.Fatalf("Unable to create server (%s), %v", *server_uri, err)
	}

	handler := www.ExampleHandler()

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	log.Printf("Listening on %s", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}
}
