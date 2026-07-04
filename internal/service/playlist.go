package service

import (
	"context"
	"database/sql"
	"fmt"

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
	playlist = lo.Ternary(playlist != nil, playlist, []repository.Playlist{})
	return playlist, err
}

func (s *PlaylistService) FindPlaylistById(ctx context.Context, id int64) (repository.Playlist, error) {
	return s.repo.FindPlaylistById(ctx, id)
}

func (s *PlaylistService) SavePlaylistAndChannels(ctx context.Context, input repository.SavePlaylistParams) (*PlaylistResult, error) {
	if input.Uri == "" {
		return nil, fmt.Errorf("url is required")
	}
	if input.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	m3uPlaylist, err := ParseM3U(ctx, input.Uri)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := s.repo.WithTx(tx)

	insertedPlayList, err := qtx.SavePlaylist(ctx, repository.SavePlaylistParams{
		Name: input.Name,
		Uri:  input.Uri,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to insert playlist. sql error: %w", err)
	}

	var tracks []repository.Channel
	for _, track := range m3uPlaylist.Channels {
		insertedChannel, err := qtx.SaveChannel(ctx, repository.SaveChannelParams{
			Name:       track.Name,
			Uri:        track.Uri,
			Length:     &track.Length,
			TvgID:      &track.TvgID,
			TvgName:    &track.TvgName,
			TvgLogo:    &track.TvgLogo,
			GroupTitle: &track.GroupTitle,
			PlaylistID: &insertedPlayList.ID,
		})
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, insertedChannel)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("unable to commit transaction: %w", err)
	}

	return &PlaylistResult{Playlist: insertedPlayList, Channels: tracks}, nil
}
