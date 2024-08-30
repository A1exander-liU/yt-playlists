package api

import (
	"context"

	playlistitems "github.com/A1exander-liU/yt-playlists/api/playlistItems"
	"github.com/A1exander-liU/yt-playlists/api/playlists"
	"github.com/A1exander-liU/yt-playlists/auth"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type ApiService struct {
	yt            *youtube.Service
	Playlists     *playlists.PlaylistService
	PlaylistItems *playlistitems.PlaylistItemService
}

func New() *ApiService {
	ctx := context.Background()
	client := auth.GetClient(ctx)
	yt, _ := youtube.NewService(ctx, option.WithHTTPClient(client))

	return &ApiService{
		yt:            yt,
		Playlists:     playlists.New(yt),
		PlaylistItems: playlistitems.New(yt),
	}
}
