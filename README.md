<h1 align="center">📁 Telegram File Proxy</h1>

[![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/deptyped/telegram-file-proxy?logo=docker)](https://hub.docker.com/r/deptyped/telegram-file-proxy)

### Introduction

This proxy allows you to provide users with links to files by `file_id` without exposing your bot's token. This is especially useful for the Telegram Web Apps feature, as it enables you to use files from Telegram directly in your web app.

To get a link to a file, simply pass `file_id` of the file as the path:

```bash
http://telegram-file-proxy/<file_id>
```

#### Query Parameters

You can also use query parameters to control response headers:

- `content-type`: Sets the `Content-Type` header. The value must be one of the
  pre-approved MIME types (e.g., `image/jpeg`, `video/mp4`, `application/pdf`).
  If the value is not allowed, the header will be omitted.

  ```bash
  http://telegram-file-proxy/<file_id>?content-type=image/png
  ```

- `filename`: Sets the `filename` in the `Content-Disposition` header. This
  suggests a name for the file when the user downloads it.

  ```bash
  http://telegram-file-proxy/<file_id>?filename=document.pdf
  ```

You can combine both parameters:

```bash
http://telegram-file-proxy/<file_id>?content-type=image/jpeg&filename=photo.jpg
```

### Usage

#### Building from source

1. Build

```bash
go mod download && go mod verify && go build -o proxy
```

2. Run

```bash
./proxy --bot-token 12345:ABCDEFGHIJKLMNOPQRSTUVWXYZ
```

💡 Pro Tip! Run `./proxy --help` to see all available command line arguments.

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

Or configuration with command line arguments:

```yaml
version: "3"
services:
  telegram-file-proxy:
    image: deptyped/telegram-file-proxy
    ports:
      - "8080:80"
    command: --bot-token <place your bot token here> --server-port 80
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

### Configuration

| ENV name            | CLI name            | Description                                                                                                            |
| ------------------- | ------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| BOT_TOKEN           | bot-token           | Bot token                                                                                                              |
| SERVER_PORT         | server-port         | Server port (8080 by default)                                                                                          |
| SERVER_HOST         | server-host         | Server hostname                                                                                                        |
| API_ROOT            | api-root            | Bot API Root (https://api.telegram.org by default)                                                                     |
| API_LOCAL           | api-local           | Allows serving files from the local filesystem when using a Local Bot API server. Set to `1` to enable.             |
| CORS_ALLOWED_ORIGIN | cors-allowed-origin | CORS allowed origin ("*" by default)                                                                                   |

Configuration is loaded in layers. Command-line flags take precedence over environment variables.
