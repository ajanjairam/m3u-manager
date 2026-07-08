package server

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/ajanjairam/m3u-manager/internal/service"
	"github.com/gofiber/fiber/v3"
)

func registerChannelRoutes(app fiber.Router, cs *service.ChannelService) {
	channelApp := app.Group("/channels")

	channelApp.Get("/", func(c fiber.Ctx) error {
		channels, err := cs.FindAllChannel(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(channels)
	})

	channelApp.Get("/page", func(c fiber.Ctx) error {
		page, err := strconv.ParseUint(c.Query("page", "1"), 10, 16)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid page number"})
		}
		pageSize, err := strconv.ParseUint(c.Query("page_size", "25"), 10, 16)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "invalid page size"})
		}
		channels, err := cs.FindAllChannelPagination(c.Context(), page, pageSize)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(channels)
	})

	channelApp.Get("/m3u", func(c fiber.Ctx) error {
		m3u, err := cs.GetM3UPlaylist(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": err.Error()})
		}
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(fiber.StatusOK).SendString(m3u)
	})

	channelApp.Get("/groups", func(c fiber.Ctx) error {
		groups, err := cs.FindAllGroup(c.Context())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.Status(404).JSON(fiber.Map{"message": "groups not found"})
			}
			return c.Status(500).JSON(fiber.Map{"message": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(groups)
	})

	channelApp.Get("/:id", func(c fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"message": "id should be an integer"})
		}
		channel, err := cs.FindById(c.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.Status(404).JSON(fiber.Map{"message": fmt.Sprintf("channel id '%d' not found", id)})
			}
			return c.Status(500).JSON(
				fiber.Map{"message": fmt.Sprintf("unable to fetch channels. sql error:%s", err.Error())})
		}
		return c.Status(fiber.StatusOK).JSON(channel)
	})
}
