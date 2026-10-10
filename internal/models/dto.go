package models

import "github.com/ajanjairam/m3u-manager/internal/repository"

type PlaylistResponse struct {
	repository.Playlist
	Channels []repository.Channel `json:"channels"`
}

type GroupResponse struct {
	repository.Playlist
	Groups []repository.ChannelGroup `json:"groups"`
}

type ChannelResponse struct {
	repository.Channel
	Group repository.ChannelGroup `json:"group"`
}

type AddFilterRequest struct {
	PlaylistID *uint    `json:"playlist"`
	GroupID    *uint    `json:"group"`
	Include    []string `json:"include"`
	Exclude    []string `json:"exclude"`
}
