# yuzu-chat 🍊

Minimal multi-provider CLI chat client written in Go.

## Quick Start

### Option 1: Git Clone (Recommended)
```bash
git clone https://github.com/icedeyes12/yuzuchat.git
cd yuzuchat
go run be.go
```

Option 2: Direct Download

```bash
# Download single file
wget https://github.com/icedeyes12/yuzuchat/raw/main/be.go

# Run directly
go run be.go
```

Option 3: Build Binary

```bash
git clone https://github.com/icedeyes12/yuzuchat.git
cd yuzuchat
go build -o yuzuchat be.go
./yuzuchat
```

Setup

1. Set Up API Keys

Create API key files in the same directory:

For Chutes AI:

```bash
echo "your_chutes_api_key_here" > cu.key
```

For OpenRouter:

```bash
echo "your_openrouter_api_key_here" > or.key
```

For Cerebras:

```bash
echo "your_cerebras_api_key_here" > ce.key
```

2. Set System Prompt (Optional)

```bash
# Edit system.txt with your rules
nano system.txt
# Or use: vim, code, etc.
```

File Structure

```
yuzuchat/
├── be.go              # Main program
├── cu.key            # Chutes AI API key
├── or.key            # OpenRouter API key  
├── ce.key            # Cerebras API key
├── system.txt        # System prompt (optional)
├── profile.json      # Settings (auto-created)
└── chat_history.json # Conversation history (auto-created)
```

Usage

Command Purpose
/key <provider> <key> Store API key
/removekey <provider> Delete API key
/provider <name> Switch provider
/model <name> Switch model
/system <text …> Set new system prompt
/system show View current prompt
/system reload Reload from disk
/models List available models
/providers List enabled providers
/clear Clear screen
/clearhistory Wipe chat history
/stream Toggle streaming mode
/info Show status
/help Show all commands
/exit or /bye Quit

Supported Providers
 
· OpenRouter 
· Chutes AI
· Cerebras 

Tips

· Edit system.txt directly for multi-line prompts
· Use /system reload after editing system.txt
· API keys are stored in separate files for security
· Conversation history keeps last 20 messages

Requirements

· Go 1.16+
· Internet connection
· API key from at least one provider

Authors

· Bani Baskara
· Yuzuki Aihara
· [guthib.com/icedeyes12/](github.com/icedeyes12/)

---

Happy chatting! 🍊💕
