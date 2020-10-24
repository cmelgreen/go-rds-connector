package main 

import (
	"context"
	"net/http"
	"os"
)

// Create router and environment then serve
func main() {
	ctx := context.Background()

	s := newServer(ctx)

	s.mux.GET("/", s.index())
	
	port := ":8050"
	if os.Getenv("PRODUCTION") == "true" {
		port = ":80"
	}

    s.log.Fatal(http.ListenAndServe(port, s.mux))
}
