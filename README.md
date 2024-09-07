<h1 align="center">yt-playlists</h1>

Simple TUI app for managing your YouTube playlists.

- [Installation](#installation)
- [Uninstall](#uninstall)
- [Setup](#setup)

## Installation

Linux

```sh
curl -s https://raw.githubusercontent.com/A1exander-liU/yt-playlists/main/scripts/install.sh | sh
```

## Uninstall

You can run the similar uninstall script.

Linux

```sh
curl -s https://raw.githubusercontent.com/A1exander-liU/yt-playlists/main/scripts/uninstall.sh | sh
```

## Setup

The app uses Google Cloud and is not production yet, a new project is required to made.

1. Go to the [Google Cloud Dev Console](https://console.cloud.google.com) and create an new project

2. In the new project, go to the menu and select "APIs & Services" > "Enabled APIs & Services".

- Search for "YouTube Data API v3" and enable it

3. Go back to the menu and select "APIs & Services" > "OAuth consent screen", select "External" and create it.

- The app name can be anything, provide your own email for user support and developer contact
- When prompted to add and remove scopes, add these 2 scopes: ".../auth/youtube.readonly" and ".../auth/youtube"
- Under test users, add your email here

4. After configuring the OAuth screen, go back to the menu "APIs & Services" > "Credentials" and create new credentials.

- Choose "OAuth Client ID"
- The name can be anything, choose "Desktop App" as the application type
- Download the json file from the created client, rename this to `client_credentials.json`

5. Move the `client_credentials.json` under:

- Linux: ~/.config/yt-playlists/
