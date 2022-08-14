<h1 align="center">📁 Telegram File Proxy</h1>

[![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/deptyped/telegram-file-proxy?logo=docker)](https://hub.docker.com/r/deptyped/telegram-file-proxy)


### Introduction

Using this proxy you can provide files to users by `file_id` without exposing the bot's token. Extremely useful for the WebApp feature to use files from Telegram in your web app.

To get a link to a file, simply pass `file_id` of the file as the path:

```bash
http://telegram-file-proxy/<file_id>
```

### Usage

#### Building from source

1. Build

```bash
go mod download && go mod verify && go build -o proxy
```

2. Set environment variables

3. Run

```bash
./proxy
```

#### Using Docker Compose

```yaml
version: "3"
services:
  telegram-file-proxy:
    image: deptyped/telegram-file-proxy
    ports:
      - "8080:80"
    environment:
      - BOT_TOKEN= # <-- place your bot token here
      - SERVER_PORT=80
```

#### Using Docker Compose with a Local Bot API Server

```yaml
version: "3"
services:
  telegram-file-proxy:
    image: deptyped/telegram-file-proxy
    ports:
      - "8080:80"
    volumes:
      - "./data:/var/lib/telegram-bot-api"
    environment:
      - BOT_TOKEN= # <-- place your bot token here
      - SERVER_PORT=80
      - API_ROOT=http://bot-api:8081
      - API_LOCAL=1

  bot-api:
    image: aiogram/telegram-bot-api:latest
    ports:
      - "8081:8081"
    volumes:
      - "./data:/var/lib/telegram-bot-api"
    environment:
      - TELEGRAM_LOCAL=1
      # Create an application with api id and api hash (get them from https://my.telegram.org/apps)
      - TELEGRAM_API_ID= # <-- place your api id here
      - TELEGRAM_API_HASH= # <-- place your api hash here
```

### Environment variables

| Variable    | Description                                                                                                            |
| ----------- | ---------------------------------------------------------------------------------------------------------------------- |
| BOT_TOKEN   | Token                                                                                                                  |
| SERVER_PORT | Server port (8080 by default)                                                                                          |
| SERVER_HOST | Server hostname                                                                                                        |
| API_ROOT    | Bot API Root (https://api.telegram.org by default)                                                                     |
| API_LOCAL   | Allow providing files from the file system, useful when using a Local Bot API with the `--local` option (0 by default) |
