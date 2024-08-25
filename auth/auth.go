package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/demo/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

func GetClient(ctx context.Context) *http.Client {
	contents, err := os.ReadFile("./client_credentials.json")
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
