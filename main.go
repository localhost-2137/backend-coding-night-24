package main

import (
	"database/sql"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	initAssistant()
	db = initDB()
	runMigrations(db)
	if err := initNasaData(); err != nil {
		log.Fatalf("Failed to fetch NASA data: %v", err)
	}

	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/assistant-ws", websocket.New(assistantWsHandler))
	app.Get("/stats", statsEndpoint)
	app.Get("/map", mapGetEndpoint)
	app.Get("/time", func(ctx *fiber.Ctx) error {
		return ctx.SendString(strconv.FormatInt(time.Now().UnixMilli(), 10))
	})

	app.Put("/map", mapPutEndpoint)

	log.Fatal(app.Listen(":3000"))
}
