# yuzu-chat

Minimal multi-provider CLI chat client written in Go.

## Providers
- Chutes AI
- OpenRouter
- Cerebras

## Quick start
1. Clone the repo.
2. Place your API key in the corresponding `*.key` file (one per line).
3. `go run be.go`

## Usage
| Command | Purpose |
|---------|---------|
| `/key <provider> <key>` | store API key |
| `/removekey <provider>` | delete API key |
| `/provider <name>` | switch provider |
| `/model <name>` | switch model |
| `/system <text …>` | set new system prompt |
| `/system show` | view current prompt |
| `/system reload` | reload from disk |
| `/models` | list available models |
| `/providers` | list enabled providers |
| `/clear` | clear screen |
| `/clearhistory` | wipe chat history |
| `/stream` | toggle streaming mode |
| `/info` | show status |
| `/exit` or `/bye` | quit |

## Files
- `system.txt` – system prompt (created on first `/system` call)
- `profile.json` – remembers last provider & model
- `chat_history.json` – last 20 exchanges

## Authors
Bani Baskara  
Yuzuki Aihara  
[guthib.com/icedeyes12/](https://github.com/icedeyes12/)
