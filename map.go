package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type pointDto struct {
	Id        int     `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Label     string  `json:"label"`
	CreatedAt string  `json:"created_at"`
}

func mapPutEndpoint(ctx *fiber.Ctx) error {
	var pointDto pointDto
	if err := ctx.BodyParser(&pointDto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	_, err := db.Exec("INSERT INTO mars_map_point (latitude, longitude, label) VALUES (?, ?, ?)",
		pointDto.Latitude, pointDto.Longitude, pointDto.Label)

	if err != nil {
		log.Errorf("Failed to add point to the map: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add point to the map",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Point added to the map",
	})
}

func mapGetEndpoint(ctx *fiber.Ctx) error {
	rows, err := db.Query("SELECT id, latitude, longitude, label, created_at FROM mars_map_point")
	if err != nil {
		log.Errorf("Failed to fetch map points: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch map points",
		})
	}

	var points []pointDto
	for rows.Next() {
		var point pointDto
		if err := rows.Scan(&point.Id, &point.Latitude, &point.Longitude, &point.Label, &point.CreatedAt); err != nil {
			log.Errorf("Failed to scan map point: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch map points",
			})
		}
		points = append(points, point)
	}

	return ctx.JSON(fiber.Map{
		"points": points,
	})
}
