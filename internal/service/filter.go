package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/ajanjairam/m3u-manager/internal/models"
	"github.com/ajanjairam/m3u-manager/internal/repository"
	"github.com/samber/lo"
)

type FilterService struct {
	db   *sql.DB
	repo *repository.Queries
}

func NewFilterService(db *sql.DB, repo *repository.Queries) *FilterService {
	return &FilterService{db: db, repo: repo}
}

func (s *FilterService) FindAllFilter(ctx context.Context) ([]models.FilterResponse, error) {
	result, err := s.repo.FindAllFilter(ctx)
	if err != nil || result == nil {
		return []models.FilterResponse{}, err
	}
	var filters []models.FilterResponse
	lo.ForEach(result, func(row repository.FindAllFilterRow, _ int) {
		filters = append(filters, models.FilterResponse{
			Filter:   row.Filter,
			Group:    row.Cgroup,
			Playlist: row.Playlist,
		})
	})
	return filters, nil
}

func (s *FilterService) SaveFilter(ctx context.Context, request models.AddFilterRequest) ([]repository.Channel, error) {
	var channels []repository.Channel
	if request.PlaylistID == nil {
		return []repository.Channel{}, errors.New("playlist is required")
	} else if (request.Include == nil || len(request.Include) == 0) && (request.Exclude == nil || len(request.Exclude) == 0) {
		return []repository.Channel{}, nil
	} else if request.GroupID == nil {
		result, err := s.repo.FindChannelByPlaylist(ctx, new(int64(*request.PlaylistID)))
		if err != nil {
			return []repository.Channel{}, err
		}
		channels = result
	} else {
		result, err := s.repo.FindChannelByPlaylistAndGroup(ctx, repository.FindChannelByPlaylistAndGroupParams{
			PlaylistID: new(int64(*request.PlaylistID)),
			GroupID:    new(int64(*request.GroupID)),
		})
		if err != nil {
			return []repository.Channel{}, err
		}
		channels = result
	}

	filteredChannels := channels
	if request.Include != nil && len(request.Include) > 0 {
		filteredChannels = lo.Filter(channels, func(channel repository.Channel, index int) bool {
			for i := range request.Include {
				contains := !strings.Contains(channel.Name, request.Include[i])
				if contains {
					return true
				}
			}
			return false
		})
	}
	if request.Exclude != nil && len(request.Exclude) > 0 {
		filteredChannels = lo.Filter(filteredChannels, func(channel repository.Channel, index int) bool {
			for i := range request.Exclude {
				contains := strings.Contains(channel.Name, request.Exclude[i])
				if contains {
					return true
				}
			}
			return false
		})
	}
	updateIDs := lo.Map(filteredChannels, func(item repository.Channel, _ int) int64 {
		return item.ID
	})
	disable, err := s.repo.UpdateChannelsDisable(ctx, updateIDs)
	if err != nil {
		return nil, err
	}
	_, err = s.repo.SaveFilter(ctx, repository.SaveFilterParams{
		Name:       "",
		PlaylistID: new(int64(*request.PlaylistID)),
		GroupID:    new(int64(*request.GroupID)),
		Include:    new(strings.Join(request.Include, ",")),
		Exclude:    new(strings.Join(request.Exclude, ",")),
	})
	if err != nil {
		return nil, err
	}
	return disable, nil
}
