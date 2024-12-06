package main

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	db = initDB()
	runMigrations(db)
	if err := initNasaData(); err != nil {
		log.Fatalf("Failed to fetch NASA data: %v", err)
	}

	app := fiber.New()

	app.Get("/stats", statsEndpoint)
	app.Get("/map", mapGetEndpoint)

	app.Put("/map", mapPutEndpoint)

	log.Fatal(app.Listen(":3000"))
}
