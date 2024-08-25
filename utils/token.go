package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

const (
	TOKEN_DIR  = ".credentials"
	TOKEN_FILE = "token.json"
)

func SaveToken(token *oauth2.Token) error {
	os.MkdirAll(filepath.Join(".", TOKEN_DIR), 0700)

	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("Failed to serialize token to json: %v", err)
	}

	os.WriteFile(filepath.Join(".", TOKEN_DIR, TOKEN_FILE), data, 0600)

	return nil
}

func LoadToken() *oauth2.Token {
	data, err := os.ReadFile(filepath.Join(".", TOKEN_DIR, TOKEN_FILE))
	if err != nil {
		return nil
	}

	var token *oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil
	}

	return token
}
