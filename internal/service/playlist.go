package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ajanjairam/m3u-manager/internal/models"
	"github.com/ajanjairam/m3u-manager/internal/repository"
	"github.com/samber/lo"
)

type PlaylistService struct {
	db   *sql.DB
	repo *repository.Queries
}

func NewPlaylistService(db *sql.DB, repo *repository.Queries) *PlaylistService {
	return &PlaylistService{db: db, repo: repo}
}

type PlaylistResult struct {
	Playlist repository.Playlist
	Channels []repository.Channel
}

func (s *PlaylistService) FindAllPlaylist(ctx context.Context) ([]repository.Playlist, error) {
	playlist, err := s.repo.FindAllPlaylist(ctx)
	return lo.Ternary(playlist != nil, playlist, []repository.Playlist{}), err
}

func (s *PlaylistService) FindPlaylistById(ctx context.Context, id int64) (repository.Playlist, error) {
	return s.repo.FindPlaylistById(ctx, id)
}

func (s *PlaylistService) SavePlaylistAndChannels(ctx context.Context, input repository.SavePlaylistParams) (models.PlaylistResponse, error) {
	if input.Uri == "" {
		return models.PlaylistResponse{}, fmt.Errorf("url is required")
	}
	if input.Name == "" {
		return models.PlaylistResponse{}, fmt.Errorf("name is required")
	}

	m3uPlaylist, err := ParseM3U(ctx, input.Uri)
	if err != nil {
		return models.PlaylistResponse{}, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.PlaylistResponse{}, fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := s.repo.WithTx(tx)

	insertedPlayList, err := qtx.SavePlaylist(ctx, repository.SavePlaylistParams{
		Name: input.Name,
		Uri:  input.Uri,
	})
	if err != nil {
		return models.PlaylistResponse{}, fmt.Errorf("unable to insert playlist. sql error: %w", err)
	}

	var channels []repository.Channel
	var groups []repository.Cgroup
	for _, track := range m3uPlaylist.Channels {
		var groupId int64
		group, isExists := lo.Find(groups, func(item repository.Cgroup) bool {
			return item.Title == track.GroupTitle
		})
		if isExists {
			groupId = group.ID
		} else {
			insertedGroup, err := qtx.SaveGroup(ctx, repository.SaveGroupParams{
				Title:      track.GroupTitle,
				PlaylistID: &insertedPlayList.ID,
			})
			if err != nil {
				return models.PlaylistResponse{}, fmt.Errorf("unable to insert group. sql error: %w", err)
			}
			groups = append(groups, insertedGroup)
			groupId = insertedGroup.ID
		}
		insertedChannel, err := qtx.SaveChannel(ctx, repository.SaveChannelParams{
			Name:       track.Name,
			Uri:        track.Uri,
			Length:     &track.Length,
			TvgID:      &track.TvgID,
			TvgName:    &track.TvgName,
			TvgLogo:    &track.TvgLogo,
			GroupID:    &groupId,
			PlaylistID: &insertedPlayList.ID,
		})
		if err != nil {
			return models.PlaylistResponse{}, err
		}
		channels = append(channels, insertedChannel)
	}

	if err := tx.Commit(); err != nil {
		return models.PlaylistResponse{}, fmt.Errorf("unable to commit transaction: %w", err)
	}

	return models.PlaylistResponse{Playlist: insertedPlayList, Channels: channels}, nil
}
