package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"

	"github.com/ajanjairam/m3u-manager/internal/models"
	"github.com/ajanjairam/m3u-manager/internal/repository"
	"github.com/jamesnetherton/m3u"
	"github.com/samber/lo"
)

type ChannelService struct {
	db   *sql.DB
	repo *repository.Queries
}

func NewChannelService(db *sql.DB, repo *repository.Queries) *ChannelService {
	return &ChannelService{db: db, repo: repo}
}

func (s *ChannelService) FindAllChannel(ctx context.Context) ([]repository.Channel, error) {
	tracks, err := s.repo.FindAllChannel(ctx)
	tracks = lo.Ternary(tracks != nil, tracks, []repository.Channel{})
	return tracks, err
}

func (s *ChannelService) FindAllChannelPagination(ctx context.Context, page uint64, pageSize uint64) (models.Pagination[repository.Channel], error) {
	tracks, err := s.repo.FindAllChannel(ctx)
	if err != nil {
		return models.NewPagination([]repository.Channel{}, 0, 0, 0), errors.New(fmt.Sprintf("unable to fetch tracks. sql error:%s", err.Error()))
	}
	tracks = lo.Ternary(tracks != nil, tracks, []repository.Channel{})
	return models.NewPagination(lo.Slice(tracks, int(((page)-1)*(pageSize)), int((page)*(pageSize))), int64(len(tracks)), int64(page), int64(pageSize)), nil
}

func (s *ChannelService) FindById(ctx context.Context, id int64) (repository.Channel, error) {
	return s.repo.FindChannelById(ctx, id)
}

func (s *ChannelService) GetM3UPlaylist(ctx context.Context) (io.Reader, error) {
	trackList, err := s.repo.FindAllActiveChannel(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get playlist. sql error: %w", err)
	}
	tracks := []m3u.Track{}
	for _, track := range trackList {
		tags := []m3u.Tag{}
		if track.TvgID != nil {
			tags = append(tags, m3u.Tag{Name: "tvg-id", Value: *track.TvgID})
		}
		if track.TvgName != nil {
			tags = append(tags, m3u.Tag{Name: "tvg-name", Value: *track.TvgName})
		}
		if track.TvgLogo != nil {
			tags = append(tags, m3u.Tag{Name: "tvg-logo", Value: *track.TvgLogo})
		}
		if track.GroupTitle != nil {
			tags = append(tags, m3u.Tag{Name: "group-title", Value: *track.GroupTitle})
		}

		tempChannel := m3u.Track{
			Name:   track.Name,
			Length: int(*track.Length),
			URI:    track.Uri,
			Tags:   tags,
		}
		tracks = append(tracks, tempChannel)
	}
	return m3u.Marshall(m3u.Playlist{
		Tracks: tracks,
	})
}
