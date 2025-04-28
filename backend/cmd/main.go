package main

import (
	"flag"
	"log"
	"os"

	"github.com/CDavidSV/Pixio/cmd/api"
	"github.com/CDavidSV/Pixio/data"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := flag.String("addr", ":3000", "HTTP network address")
	flag.Parse()

	dsn := os.Getenv("DATABASE_URL")

	db, err := data.NewPostgresPool(dsn)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	server := api.NewServer(*addr, db)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
