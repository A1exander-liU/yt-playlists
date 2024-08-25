package playlistitems

import (
	"context"
	"fmt"

	"google.golang.org/api/youtube/v3"
)

type PlaylistItemService struct {
	yt *youtube.Service
}

func New(yt *youtube.Service) *PlaylistItemService {
	return &PlaylistItemService{yt: yt}
}

func list(yt *youtube.Service, playlistId string, part []string) ([]*youtube.PlaylistItem, error) {
	items := []*youtube.PlaylistItem{}
	req := yt.PlaylistItems.List(part).PlaylistId(playlistId).MaxResults(50)
	err := req.Pages(context.Background(), func(res *youtube.PlaylistItemListResponse) error {
		items = append(items, res.Items...)
		return nil
	})

	return items, err
}

func addToPlaylist(yt *youtube.Service, playlistId string, playlistItem *youtube.PlaylistItem) (*youtube.PlaylistItem, error) {
	playlistItemToAdd := &youtube.PlaylistItem{Snippet: &youtube.PlaylistItemSnippet{PlaylistId: playlistId, ResourceId: playlistItem.Snippet.ResourceId}}
	req := yt.PlaylistItems.Insert([]string{"snippet"}, playlistItemToAdd)

	res, err := req.Do()
	return res, err
}

func deleteFromPlaylist(yt *youtube.Service, id string) error {
	req := yt.PlaylistItems.Delete(id)

	err := req.Do()
	return err
}

func moveToPlaylist(yt *youtube.Service, playlistId string, playlistItem *youtube.PlaylistItem) (*youtube.PlaylistItem, error) {
	deleteErr := deleteFromPlaylist(yt, playlistItem.Id)
	if deleteErr != nil {
		return nil, deleteErr
	}

	res, addErr := addToPlaylist(yt, playlistId, playlistItem)
	return res, addErr
}

func (playlistItemService *PlaylistItemService) List(playlistId string, part []string) ([]*youtube.PlaylistItem, error) {
	res, err := list(playlistItemService.yt, playlistId, part)
	return res, err
}

func (playlistItemService *PlaylistItemService) Add(playlistId string, playlistItems []*youtube.PlaylistItem) ([]*youtube.PlaylistItem, error) {
	added := []*youtube.PlaylistItem{}
	errs := []error{}

	for _, playlistItem := range playlistItems {
		res, err := addToPlaylist(playlistItemService.yt, playlistId, playlistItem)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		added = append(added, res)
	}

	if len(errs) > 0 {
		return added, fmt.Errorf("Errors occured while adding: %v", errs)
	}
	return added, nil
}

func (playlistItemService *PlaylistItemService) Delete(ids []string) ([]string, error) {
	deleted := []string{}
	errs := []error{}

	for _, id := range ids {
		err := deleteFromPlaylist(playlistItemService.yt, id)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		deleted = append(deleted, id)
	}

	if len(errs) > 0 {
		return deleted, fmt.Errorf("Errors occured while deleting: %v", errs)
	}
	return deleted, nil
}

func (playlistItemService *PlaylistItemService) Move(playlistId string, playlistItems []*youtube.PlaylistItem) ([]*youtube.PlaylistItem, error) {
	moved := []*youtube.PlaylistItem{}
	errs := []error{}

	for _, playlistItem := range playlistItems {
		res, err := moveToPlaylist(playlistItemService.yt, playlistId, playlistItem)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		moved = append(moved, res)
	}

	if len(errs) > 0 {
		return moved, fmt.Errorf("Errors occured while moving: %v", errs)
	}
	return moved, nil
}

func (playlistItemService *PlaylistItemService) Exists(playlistId string, playlistItem *youtube.PlaylistItem) (bool, error) {
	videos, err := list(playlistItemService.yt, playlistId, []string{"snippet"})
	if err != nil {
		return false, err
	}

	for _, video := range videos {
		if playlistItem.Snippet.ResourceId.VideoId == video.Snippet.ResourceId.VideoId {
			return true, nil
		}
	}

	return false, nil
}
