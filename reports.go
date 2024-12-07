package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type reportDto struct {
	Id        int    `json:"id"`
	Label     string `json:"label"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func getReportsEndpoint(ctx *fiber.Ctx) error {
	rows, err := db.Query("SELECT id, label, content, created_at FROM reports ORDER BY created_at DESC")
	if err != nil {
		log.Errorf("Failed to fetch reports: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch reports",
		})
	}
	defer rows.Close()

	var reports []reportDto
	for rows.Next() {
		var report reportDto
		if err := rows.Scan(&report.Id, &report.Label, &report.Content, &report.CreatedAt); err != nil {
			log.Errorf("Failed to scan report: %v", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch reports",
			})
		}
		reports = append(reports, report)
	}

	return ctx.JSON(reports)
}
