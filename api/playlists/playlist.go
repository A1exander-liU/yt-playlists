package playlists

import (
	"fmt"

	"google.golang.org/api/youtube/v3"
)

type PlaylistService struct {
	yt *youtube.Service
}

func New(yt *youtube.Service) *PlaylistService {
	return &PlaylistService{yt: yt}
}

func listPlaylists(yt *youtube.Service, part []string, nextPage string) (*youtube.PlaylistListResponse, error) {
	req := yt.Playlists.List(part).Mine(true).PageToken(nextPage).MaxResults(50)

	res, err := req.Do()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (playlistService *PlaylistService) List(part []string) (*youtube.PlaylistListResponse, error) {
	res, err := listPlaylists(playlistService.yt, part, "")
	playlists := []*youtube.Playlist{}
	if err != nil {
		return nil, err
	}
	playlists = append(playlists, res.Items...)

	for res.NextPageToken != "" {
		res, err = listPlaylists(playlistService.yt, part, res.NextPageToken)
		if err != nil {
			return nil, err
		}
		playlists = append(playlists, res.Items...)
	}

	res.Items = playlists
	return res, nil
}

func insertPlaylist(yt *youtube.Service, name string, description string, status string) (*youtube.Playlist, error) {
	if status != "public" && status != "private" && status != "unlisted" {
		return nil, fmt.Errorf("%v is not a valid staus: Must be one of 'private', 'public', or 'unlisted'", status)
	}
	newPlaylistStatus := &youtube.PlaylistStatus{PrivacyStatus: status}
	newPlaylist := &youtube.Playlist{Snippet: &youtube.PlaylistSnippet{Title: name, Description: description}, Status: newPlaylistStatus}
	req := yt.Playlists.Insert([]string{"snippet,status"}, newPlaylist)

	res, err := req.Do()
	return res, err
}

func (playlistService *PlaylistService) Insert(name string, description string, status string) (*youtube.Playlist, error) {
	res, err := insertPlaylist(playlistService.yt, name, description, status)
	return res, err
}

func deletePlaylist(yt *youtube.Service, playlistId string) error {
	req := yt.Playlists.Delete(playlistId)

	err := req.Do()
	return err
}

func (playlistService *PlaylistService) Delete(playlistId string) error {
	err := deletePlaylist(playlistService.yt, playlistId)
	return err
}
