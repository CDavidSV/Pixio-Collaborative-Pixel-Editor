package config

import (
	"log"
	"os"
	"time"

	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %s", err)
	}

	DSN = os.Getenv("DATABASE_URL")
	AccessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
	RefreshTokenSecret = os.Getenv("REFRESH_TOKEN_SECRET")
}

var (
	PixioLogo = `
    ____     _             _                   ___     ____     ____
   / __ \   (_)   _  __   (_)  ____           /   |   / __ \   /  _/
  / /_/ /  / /   | |/_/  / /  / __ \         / /| |  / /_/ /   / /
 / ____/  / /   _>  <   / /  / /_/ /        / ___ | / ____/  _/ /
/_/      /_/   /_/|_|  /_/   \____/        /_/  |_|/_/      /___/
`
	DSN                string
	AccessTokenSecret  string
	RefreshTokenSecret string
	CorsConfig         = cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "HEAD", "OPTION", "PUT"},
		AllowedHeaders:   []string{"User-Agent", "Content-Type", "Accept", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "DNT", "Host", "Origin", "Pragma", "Referer", "Cookie"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}

	SessionExpiration     = time.Hour * 24 * 30 // 30 days
	AccessTokenExpiration = time.Minute * 15    // 15 minutes
)
