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
	channels, err := s.repo.FindAllChannel(ctx)
	return lo.Ternary(channels != nil, channels, []repository.Channel{}), err
}

func (s *ChannelService) FindAllChannelPagination(ctx context.Context, page uint64, pageSize uint64) (models.Pagination[models.ChannelResponse], error) {
	result, err := s.repo.FindAllChannelAndGroup(ctx)
	if err != nil {
		return models.NewPagination([]models.ChannelResponse{}, 0, 0, 0), errors.New(fmt.Sprintf("unable to fetch channels. sql error:%s", err.Error()))
	}
	var channels []models.ChannelResponse
	lo.ForEach(result, func(row repository.FindAllChannelAndGroupRow, index int) {
		channels = append(channels, models.ChannelResponse{
			Channel: row.Channel,
			Group:   row.Cgroup,
		})
	})
	channels = lo.Ternary(channels != nil, channels, []models.ChannelResponse{})
	return models.NewPagination(lo.Slice(channels, int((page-1)*pageSize), int(page*pageSize)), int64(len(channels)), int64(page), int64(pageSize)), nil
}

func (s *ChannelService) FindById(ctx context.Context, id int64) (repository.Channel, error) {
	return s.repo.FindChannelById(ctx, id)
}

func (s *ChannelService) GetM3UPlaylist(ctx context.Context) (string, error) {
	result, err := s.repo.FindAllActiveChannel(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to get playlist. sql error: %w", err)
	}
	var channels []models.Channel
	lo.ForEach(result, func(row repository.FindAllActiveChannelRow, _ int) {
		channels = append(channels, models.Channel{
			Length:     *row.Channel.Length,
			TvgID:      *row.Channel.TvgID,
			TvgName:    *row.Channel.TvgName,
			TvgLogo:    *row.Channel.TvgLogo,
			GroupTitle: row.Cgroup.Title,
			Name:       row.Channel.Name,
			Uri:        row.Channel.Uri,
		})
	})
	return MarshalM3U(channels)
}
