package service

import (
	"context"

	"github.com/ajanjairam/m3u-manager/internal/models"
	"github.com/ajanjairam/m3u-manager/internal/repository"
	"github.com/samber/lo"
)

func (s *ChannelService) FindAllGroup(ctx context.Context) ([]models.GroupResponse, error) {
	groupsAndPlaylists, err := s.repo.FindAllGroup(ctx)
	if err != nil {
		return nil, err
	}
	var playlist []models.GroupResponse
	lo.ForEach(groupsAndPlaylists, func(row repository.FindAllGroupRow, _ int) {
		_, index, isExist := lo.FindIndexOf(playlist, func(item models.GroupResponse) bool {
			return item.ID == row.Playlist.ID
		})
		if isExist {
			playlist[index].Groups = lo.UniqBy(append(playlist[index].Groups, row.Cgroup), func(item repository.Cgroup) string {
				return item.Title
			})
		} else {
			playlist = append(playlist, models.GroupResponse{
				Playlist: row.Playlist,
				Groups:   []repository.Cgroup{row.Cgroup},
			})
		}
	})
	return playlist, nil
}
