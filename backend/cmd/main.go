package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/CDavidSV/Pixio/cmd/api"
	"github.com/CDavidSV/Pixio/config"
	"github.com/CDavidSV/Pixio/data"
)

func main() {
	fmt.Println(config.PixioLogo)
	fmt.Println("Pixio API - Version 1.0.0")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	addr := flag.String("addr", ":3000", "HTTP network address")
	flag.Parse()

	db, err := data.NewPostgresPool(config.DSN)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	logger.Info("Connected to the database")
	logger.Info("Starting server on", "address", *addr)
	server := api.NewServer(*addr, db)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
