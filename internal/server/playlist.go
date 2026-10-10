package server

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/ajanjairam/m3u-manager/internal/repository"
	"github.com/ajanjairam/m3u-manager/internal/service"
	"github.com/gofiber/fiber/v3"
)

func registerPlaylistRoutes(app fiber.Router, ps *service.PlaylistService) {

	playlistApp := app.Group("/playlist")

	playlistApp.Get("/", func(c fiber.Ctx) error {
		allPlaylist, err := ps.FindAllPlaylist(c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(allPlaylist)
	})

	playlistApp.Get("/:id", func(c fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "id should be an integer"})
		}
		playlist, err := ps.FindPlaylistById(c.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.Status(404).JSON(fiber.Map{"message": fmt.Sprintf("playlist id '%d' not found", id)})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(playlist)
	})

	playlistApp.Post("/", func(c fiber.Ctx) error {
		var body struct {
			Name string `json:"name"`
			URI  string `json:"uri"`
		}
		if err := c.Bind().All(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}

		result, err := ps.SavePlaylistAndChannels(c.Context(), repository.SavePlaylistParams{
			Name: body.Name,
			Uri:  body.URI,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"playlist": result.Playlist, "channels": result.Channels})
	})

	playlistApp.Delete("/:id", func(c fiber.Ctx) error {
		playlistId, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if(err != nil){
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		}
		channels, err := strconv.ParseBool(c.Query("channels"))
		if(err != nil){
			channels = true
		}
		var success bool
		if(channels){
			success, err = ps.DeletePlaylistAndChannels(c, playlistId)
		} else {
			success, err = ps.DeletePlaylist(c, playlistId)
		}
		if(err != nil || !success){
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		} else {
			return c.SendStatus(fiber.StatusNoContent)
		}
	})
}
