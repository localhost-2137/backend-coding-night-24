package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type itemDto struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Quantity      int    `json:"quantity"`
	Unit          string `json:"unit"`
	CriticalLevel int    `json:"critical_level"`
	LastChecked   string `json:"last_checked"`
	Location      string `json:"location"`
}

func getCurrentItems() ([]itemDto, error) {
	rows, err := db.Query("SELECT * FROM items")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []itemDto
	for rows.Next() {
		var item itemDto
		if err := rows.Scan(&item.Id, &item.Name, &item.Description, &item.Quantity, &item.Unit, &item.CriticalLevel, &item.LastChecked, &item.Location); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func GetItemsEndpoint(ctx *fiber.Ctx) error {
	rows, err := db.Query("SELECT * FROM items")

	if err != nil {
		log.Errorf("Failed to query items: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get items in storage",
		})
	}

	var items []itemDto

	for rows.Next() {
		var item itemDto
		if err := rows.Scan(&item.Id, &item.Name, &item.Description, &item.Quantity, &item.Unit, &item.CriticalLevel, &item.LastChecked, &item.Location); err != nil {
			log.Errorf("Failed to query items: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get items in storage",
			})
		}

		items = append(items, item)
	}

	return ctx.Status(fiber.StatusOK).JSON(items)
}
