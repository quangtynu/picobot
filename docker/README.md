# Docker Deployment

Run Picobot as a Docker container — one command to start.

## Quick Start

### Option 1: Docker Compose (Recommended)

```sh
# 1. Create your .env file
cp docker/.env.example docker/.env

# 2. Edit .env with your API key and settings
nano docker/.env

# 3. Start
docker compose -f docker/docker-compose.yml up -d

# 4. Check logs
docker compose -f docker/docker-compose.yml logs -f
```

### Option 2: Docker Run

```sh
# Build the image
docker build -f docker/Dockerfile -t picobot .

# Run with environment variables
docker run -d \
  --name picobot \
  --restart unless-stopped \
  -e OPENROUTER_API_KEY="sk-or-v1-YOUR_KEY" \
  -e PICOBOT_MODEL="google/gemini-2.5-flash" \
  -e TELEGRAM_BOT_TOKEN="123456:ABC..." \
  -e TELEGRAM_ALLOW_FROM="8281248569" \
  -v picobot-data:/home/picobot/.picobot \
  picobot
```

### Option 3: Mount Your Own Config

If you already have a `config.json`, mount it directly:

```sh
docker run -d \
  --name picobot \
  --restart unless-stopped \
  -v /path/to/your/config.json:/home/picobot/.picobot/config.json \
  -v picobot-data:/home/picobot/.picobot/workspace \
  picobot
```

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `OPENROUTER_API_KEY` | Yes | — | Your OpenRouter API key |
| `PICOBOT_MODEL` | No | `google/gemini-2.5-flash` | LLM model to use |
| `TELEGRAM_BOT_TOKEN` | No | — | Telegram bot token from @BotFather |
| `TELEGRAM_ALLOW_FROM` | No | — | Comma-separated Telegram user IDs |

## Management

```sh
# Stop
docker compose -f docker/docker-compose.yml down

# Rebuild after code changes
docker compose -f docker/docker-compose.yml up -d --build

# View logs
docker compose -f docker/docker-compose.yml logs -f

# Shell into container
docker exec -it picobot sh
```

## Data Persistence

All data is stored in the `picobot-data` Docker volume:
- `config.json` — configuration
- `workspace/` — bootstrap files, memory, skills

Data persists across container restarts and rebuilds.
