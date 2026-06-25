# Legacy iOS AI Telegram Bot

This repository contains a Golang-based Telegram bot designed to serve as a frontend bridge to a LiteLLM backend. 

This project is part of an architecture intended to revive older devices (like iPads running iOS 9-12) by bypassing outdated browsers that cannot run modern AI web apps. Using Telegram as a native UI and a Tailscale Zero Trust network, you can securely route messages from the iPad to a modern AI service (like Gemini or OpenAI) processed via a LiteLLM gateway.

## Architecture Diagram

`iPad (Telegram)` -> `Golang Bot (iSH or Tiny VM)` -> `Tailscale Mesh` -> `LiteLLM Backend` -> `LLM API (Gemini, etc.)`

## Setup & Deployment

### 1. Prerequisites
- **Telegram Bot Token:** Use BotFather in Telegram to create a new bot and obtain an API token.
- **LiteLLM Gateway:** Ensure you have LiteLLM running and bound to your Tailscale network IP. Obtain a virtual key.
- **Tailscale:** Both the bot host and the LiteLLM host must be on the same Tailscale WireGuard mesh for secure communication.

### 2. Environment Variables
The bot requires the following environment variables:
- `TELEGRAM_BOT_TOKEN`: The API token from BotFather.
- `LITELLM_URL`: The URL to your LiteLLM chat completions endpoint (e.g., `http://100.x.x.x:4000/v1/chat/completions`).
- `LITELLM_KEY`: Your LiteLLM virtual key.

### 3. Compilation
Because Go compiles to a single, zero-dependency binary, you can run it on almost any system, including local iPad execution via the `iSH` emulator.

```bash
# Standard compilation (e.g., for a Linux VM)
go build -o telegram-bot main.go

# Cross-compilation for iSH (32-bit x86 emulator on iPad)
GOOS=linux GOARCH=386 go build -o telegram-bot main.go
```

### 4. Running the Bot

**Option A: Local on iPad via iSH**
1. Transfer the compiled 32-bit `telegram-bot` binary to your iPad's iSH environment.
2. Ensure the Tailscale iOS app is connected.
3. Run the bot:
```sh
TELEGRAM_BOT_TOKEN="your_token" LITELLM_URL="your_url" LITELLM_KEY="your_key" ./telegram-bot
```

**Option B: On a dedicated Cloud VM**
1. Transfer the standard binary to your VM.
2. Run it inside a `tmux` session or create a `systemd` service for 24/7 uptime.

## Reference Design
For the full conceptual guide, please refer to [design.md](./design.md) included in this repository.
