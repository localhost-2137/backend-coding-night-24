package main

import (
	"github.com/gofiber/fiber/v2"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := initDB()
	runMigrations(db)
	if err := initNasaData(); err != nil {
		log.Fatalf("Failed to fetch NASA data: %v", err)
	}

	app := fiber.New()

	app.Get("/stats", statsEndpoint)

	log.Fatal(app.Listen(":3000"))
}
