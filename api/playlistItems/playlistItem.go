package playlistitems

import "google.golang.org/api/youtube/v3"

type PlaylistItemService struct {
	yt *youtube.Service
}

func New(yt *youtube.Service) *PlaylistItemService {
	return &PlaylistItemService{yt: yt}
}

func list(yt *youtube.Service, playlistId string, part []string, nextPage string) (*youtube.PlaylistItemListResponse, error) {
	req := yt.PlaylistItems.List(part).PlaylistId(playlistId).MaxResults(50).PageToken(nextPage)

	res, err := req.Do()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (playlistItemService *PlaylistItemService) List(playlistId string, part []string) (*youtube.PlaylistItemListResponse, error) {
	res, err := list(playlistItemService.yt, playlistId, part, "")
	videos := []*youtube.PlaylistItem{}
	if err != nil {
		return nil, err
	}
	videos = append(videos, res.Items...)

	for res.NextPageToken != "" {
		res, err = list(playlistItemService.yt, playlistId, part, "")
		if err != nil {
			return nil, err
		}
		videos = append(videos, res.Items...)
	}

	res.Items = videos
	return res, nil
}
