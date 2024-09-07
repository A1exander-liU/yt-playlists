#!/bin/sh

APP_NAME="yt-playlists"
OS="Linux"
ARCH=$(uname -m)
REPO_URL="https://github.com/A1exander-liU/$APP_NAME"
VERSION="v0.9.2"

case $(uname -m) in
x86_64)
  ARCH="x86_64"
  ;;
i386 | i686)
  ARCH="i386"
  ;;
aarch64)
  ARCH="arm64"
  ;;
*)
  command ...
  ;;
esac

# Make directories for the app
mkdir -p "$HOME/.config/$APP_NAME"
mkdir -p "$HOME/.$APP_NAME"

PACKAGE="${APP_NAME}_${OS}_${ARCH}.tar.gz"
DOWNLOAD="$REPO_URL/releases/download/$VERSION/$PACKAGE"

curl -o ~/Downloads/$PACKAGE -L $DOWNLOAD
tar -xzvf ~/Downloads/$PACKAGE
mv $APP_NAME /usr/bin/
