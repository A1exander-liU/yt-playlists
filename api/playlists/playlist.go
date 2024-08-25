package playlists

import (
	"context"
	"fmt"

	"google.golang.org/api/youtube/v3"
)

type PlaylistService struct {
	yt *youtube.Service
}

func New(yt *youtube.Service) *PlaylistService {
	return &PlaylistService{yt: yt}
}

func listPlaylists(yt *youtube.Service, part []string) ([]*youtube.Playlist, error) {
	items := []*youtube.Playlist{}
	req := yt.Playlists.List(part).Mine(true).MaxResults(50)

	err := req.Pages(context.Background(), func(res *youtube.PlaylistListResponse) error {
		items = append(items, res.Items...)
		return nil
	})

	return items, err
}

func insertPlaylist(yt *youtube.Service, name string, description string, status string) (*youtube.Playlist, error) {
	if status != "public" && status != "private" && status != "unlisted" {
		return nil, fmt.Errorf("%v is not a valid status: Must be one of 'private', 'public', or 'unlisted'", status)
	}
	newPlaylistStatus := &youtube.PlaylistStatus{PrivacyStatus: status}
	newPlaylist := &youtube.Playlist{Snippet: &youtube.PlaylistSnippet{Title: name, Description: description}, Status: newPlaylistStatus}
	req := yt.Playlists.Insert([]string{"snippet,status"}, newPlaylist)

	res, err := req.Do()
	return res, err
}

func deletePlaylist(yt *youtube.Service, playlistId string) error {
	req := yt.Playlists.Delete(playlistId)

	err := req.Do()
	return err
}

func (playlistService *PlaylistService) List(part []string) ([]*youtube.Playlist, error) {
	res, err := listPlaylists(playlistService.yt, part)
	return res, err
}

func (playlistService *PlaylistService) Insert(name string, description string, status string) (*youtube.Playlist, error) {
	res, err := insertPlaylist(playlistService.yt, name, description, status)
	return res, err
}

func (playlistService *PlaylistService) Delete(playlistId string) error {
	err := deletePlaylist(playlistService.yt, playlistId)
	return err
}
