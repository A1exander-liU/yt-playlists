package playlists

import "google.golang.org/api/youtube/v3"

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
