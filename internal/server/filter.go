package server

import (
	"github.com/ajanjairam/m3u-manager/internal/models"
	"github.com/ajanjairam/m3u-manager/internal/service"
	"github.com/gofiber/fiber/v3"
)

func registerFilterRoutes(app fiber.Router, fis *service.FilterService) {
	filterApp := app.Group("/filters")

	filterApp.Get("/", func(c fiber.Ctx) error {
		filters, err := fis.FindAllFilter(c.Context())
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(filters)
	})

	filterApp.Post("/", func(c fiber.Ctx) error {
		var requestBody models.AddFilterRequest
		if err := c.Bind().All(&requestBody); err != nil {
			return err
		}
		filters, err := fis.SaveFilter(c.Context(), requestBody)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(filters)
	})
}
