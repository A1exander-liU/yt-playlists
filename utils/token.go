package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/oauth2"
)

const (
	APP_NAME   = "yt-playlists"
	TOKEN_DIR  = ".credentials"
	TOKEN_FILE = "token.json"
)

func TokenPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Unable to locate home directory")
		return ""
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("LOCALAPPDATA"), APP_NAME, "Data", TOKEN_DIR)
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", APP_NAME, "Data", TOKEN_DIR)
	case "linux":
		return filepath.Join(homeDir, fmt.Sprintf(".%s", APP_NAME), TOKEN_DIR)
	default:
		return filepath.Join(homeDir, fmt.Sprintf(".%s", APP_NAME), TOKEN_DIR)
	}
}

func SaveToken(token *oauth2.Token) error {
	os.MkdirAll(TokenPath(), 0700)

	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("Failed to serialize token to json: %v", err)
	}

	os.WriteFile(filepath.Join(TokenPath(), TOKEN_FILE), data, 0600)

	return nil
}

func LoadToken() *oauth2.Token {
	tokenPath := filepath.Join(TokenPath(), TOKEN_FILE)
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return nil
	}

	var token *oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil
	}

	return token
}

func RemoveToken() {
	tokenPath := filepath.Join(TokenPath(), TOKEN_FILE)
	os.Remove(tokenPath)
}
