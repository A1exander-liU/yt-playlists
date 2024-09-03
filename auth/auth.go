package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/A1exander-liU/yt-playlists/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

const APP_CLIENT = "client_credentials.json"

func AppClientPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Unable not locate home directory")
		return ""
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("LOCALAPPDATA"), utils.APP_NAME)
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", utils.APP_NAME)
	case "linux":
		return filepath.Join(homeDir, ".config", utils.APP_NAME)
	default:
		return filepath.Join(homeDir, ".config", utils.APP_NAME)
	}
}

func GetClient(ctx context.Context) *http.Client {
	appClientPath := filepath.Join(AppClientPath(), APP_CLIENT)
	contents, err := os.ReadFile(appClientPath)
	if err != nil {
		log.Fatalf("Error occurred reading the file: %v", err)
	}

	config, err := google.ConfigFromJSON(contents, youtube.YoutubeScope, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Failed to parse credentials file: %v", err)
	}

	var token *oauth2.Token
	token = utils.LoadToken()
	if token == nil {
		token = getTokenFromWeb(config)
	}

	return config.Client(ctx, token)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Go to the link below and obtain authorization code\n%v\nEnter code: ", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatal("Unable to retrieve token: %v", err)
	}

	utils.SaveToken(tok)
	return tok
}
