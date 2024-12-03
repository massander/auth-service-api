package main

import (
	"log"
	"os"

	"auth-service-api/server"
	"auth-service-api/storage/postgres"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	// Storage
	storage, err := postgres.New(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	var migrationsFolder string
	if os.Getenv("MIGRATIONS_FOLDER") == "" {
		migrationsFolder = server.MIGRATIONS_FOLDER
	}

	if err := storage.Migrate(migrationsFolder); err != nil {
		log.Fatal(err)
	}

	// Server
	server := server.New(storage)

	server.Start()

}
