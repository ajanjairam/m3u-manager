package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ajanjairam/m3u-manager/internal/models"
	"github.com/ajanjairam/m3u-manager/internal/repository"
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

func (s *ChannelService) GetM3UPlaylist(ctx context.Context) (*string, error) {
	channelList, err := s.repo.FindAllActiveChannel(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get playlist. sql error: %w", err)
	}
	var channels []models.Channel
	for _, channel := range channelList {
		channels = append(channels, models.Channel{
			Length:     *channel.Length,
			TvgID:      *channel.TvgID,
			TvgName:    *channel.TvgName,
			TvgLogo:    *channel.TvgLogo,
			GroupTitle: *channel.GroupTitle,
			Name:       channel.Name,
			Uri:        channel.Uri,
		})
	}
	return MarshalM3U(channels)
}
