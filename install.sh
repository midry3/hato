#!/bin/sh
set -eu

DEST=""

while [ $# -gt 0 ]; do
  case "$1" in
    -d|--dir)
      DEST=$2
      shift 2
      ;;
    *)
      echo "unknown option: $1" >&2
      exit 1
      ;;
  esac
done

if [ -z "$DEST" ]; then
  if [ -d "$HOME/.local/bin" ]; then
    DEST="$HOME/.local/bin"
  elif [ -d "$HOME/bin" ]; then
    DEST="$HOME/bin"
  else
    DEST="$HOME/.local/bin"
    mkdir -p "$DEST"
  fi
else
  mkdir -p "$DEST"
fi

OS=$(uname -s)
case "$OS" in
  Linux)   OS=linux ;;
  Darwin)  OS=darwin ;;
  *)
    echo "unsupported OS: $OS" >&2
    exit 1
    ;;
esac

CPU=$(uname -m)
case "$CPU" in
  x86_64|amd64) CPU=amd64 ;;
  aarch64|arm64) CPU=arm64 ;;
  *)
    echo "unsupported CPU: $CPU" >&2
    exit 1
    ;;
esac

FILE="https://github.com/midry3/hato/releases/latest/download/hato_${OS}_${CPU}"
curl -Lo "$DEST/hato" $FILE
chmod +x "$DEST/hato"
printf "Installed to '\033[32m$DEST\033[0m'\n"
case ":$PATH:" in
  *":$DEST:"*) ;;
  *) printf "\033[31mNOTE\033[0m: You should add '$DEST' to \033[33m\$PATH\033[0m\n" ;;
esac
