package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/CDavidSV/Pixio/cmd/api/config"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	fmt.Println(config.PixioLogo)

	server := http.Server{
		Addr:         *addr,
		Handler:      loadRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server.ListenAndServe()
}
