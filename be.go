
                // be.go - YuzuChat Multi-Provider AI Client  
            // ©2025 hkkm project | built with love 💕
        // guthib.com/icedeyes12

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Color string

const (
	Red    Color = "\033[31m"
	Green  Color = "\033[32m"
	Yellow Color = "\033[33m"
	Cyan   Color = "\033[36m"
	Purple Color = "\033[35m"
	Reset  Color = "\033[0m"
)

func colorPrint(color Color, message string, args ...interface{}) {
	fmt.Printf(string(color)+message+string(Reset), args...)
}

type Message struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	Model     string `json:"model"`
	Provider  string `json:"provider"`
}

type AIProvider struct {
	Name      string
	BaseURL   string
	APIKey    string
	Models    []string
	IsEnabled bool
	KeyFile   string
}

type YuzuChat struct {
	providers           map[string]*AIProvider
	currentProvider     string
	historyFile         string
	profileFile         string
	systemFile          string
	conversationHistory []Message
	systemPrompt        string
	model               string
}

func NewYuzuChat(historyFile, profileFile, systemFile string) *YuzuChat {
	chat := &YuzuChat{
		historyFile:     historyFile,
		profileFile:     profileFile,
		systemFile:      systemFile,
		providers:       make(map[string]*AIProvider),
		currentProvider: "chutes",
		model:           "deepseek-ai/DeepSeek-V3-0324",
	}
	chat.loadProviders()
	chat.loadProfile()
	chat.loadSystemPrompt()
	chat.loadHistory()
	return chat
}

func (y *YuzuChat) loadProviders() {
	y.providers["chutes"] = &AIProvider{
		Name:    "Chutes AI",
		BaseURL: "https://llm.chutes.ai/v1/chat/completions",
		KeyFile: "cu.key",
		Models: []string{
			"deepseek-ai/DeepSeek-V3-0324",
			"deepseek-ai/DeepSeek-V3.1-Terminus",
			"tngtech/DeepSeek-R1T-Chimera",
			"tngtech/DeepSeek-R1T2-Chimera",
			"Qwen/Qwen3-235B-A22B-Instruct",
			"Qwen/Qwen3-VL-235B-A22B-Thinking",
			"Qwen/Qwen3-Coder-480B-A35B-Instruct-FP8",
			"zai-org/GLM-4.5-FP8",
			"zai-org/GLM-4.6-FP8",
			"deepseek-ai/DeepSeek-R1",
		},
	}
	y.providers["openrouter"] = &AIProvider{
		Name:    "OpenRouter",
		BaseURL: "https://openrouter.ai/api/v1/chat/completions",
		KeyFile: "or.key",
		Models: []string{
			"tngtech/deepseek-r1t2-chimera:free",
			"z_ai/glm-4.5-air:free",
			"tngtech/deepseek-r1t-chimera:free",
			"deepseek/deepseek-v3:free",
			"deepseek/r1:free",
			"qwen/qwen3-235b-a22b:free",
			"meituan/longcat-flash-chat:free",
		},
	}
	y.providers["cerebras"] = &AIProvider{
		Name:    "Cerebras",
		BaseURL: "https://api.cerebras.ai/v1/chat/completions",
		KeyFile: "ce.key",
		Models: []string{
			"qwen-3-235b-a22b-instruct-2507",
			"qwen-3-235b-a22b-thinking-2507",
			"qwen-3-coder-480b",
			"qwen-3-32b",
			"gpt-oss-120b",
			"llama-3.3-70b",
			"llama-4-scout-17b-16e-instruct",
			"llama3.1-8b",
		},
	}
	enabledCount := 0
	for name, provider := range y.providers {
		provider.APIKey = y.loadKeyFile(provider.KeyFile)
		provider.IsEnabled = provider.APIKey != ""
		if provider.IsEnabled {
			enabledCount++
			colorPrint(Green, "✅ %s: API key loaded from %s\n", name, provider.KeyFile)
		} else {
			colorPrint(Yellow, "⚠️ %s: No API key found in %s\n", name, provider.KeyFile)
		}
	}
	colorPrint(Green, "\n🎯 Total providers enabled: %d/%d\n", enabledCount, len(y.providers))
}

func (y *YuzuChat) loadKeyFile(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func (y *YuzuChat) saveKeyFile(filename, apiKey string) error {
	return os.WriteFile(filename, []byte(apiKey), 0600)
}

func (y *YuzuChat) removeKeyFile(filename string) error {
	err := os.Remove(filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (y *YuzuChat) loadProfile() {
	data, err := os.ReadFile(y.profileFile)
	if err != nil {
		if os.IsNotExist(err) {
			colorPrint(Yellow, "📝 No profile found, starting with default settings\n")
			return
		}
		colorPrint(Red, "❌ Error loading profile: %v\n", err)
		return
	}
	var profileData struct {
		Model    string `json:"model"`
		Provider string `json:"provider"`
	}
	if err := json.Unmarshal(data, &profileData); err != nil {
		colorPrint(Red, "❌ Error parsing profile: %v\n", err)
		return
	}
	if profileData.Model != "" {
		y.model = profileData.Model
	}
	if profileData.Provider != "" {
		y.currentProvider = profileData.Provider
	}
	colorPrint(Green, "📖 Profile loaded: %s provider, %s model\n", y.currentProvider, y.model)
}

func (y *YuzuChat) saveProfile() {
	profileData := struct {
		Model       string `json:"model"`
		Provider    string `json:"provider"`
		LastUpdated string `json:"last_updated"`
	}{
		Model:       y.model,
		Provider:    y.currentProvider,
		LastUpdated: time.Now().Format(time.RFC3339),
	}
	data, err := json.MarshalIndent(profileData, "", "  ")
	if err != nil {
		colorPrint(Red, "❌ Error marshaling profile: %v\n", err)
		return
	}
	if err := os.WriteFile(y.profileFile, data, 0644); err != nil {
		colorPrint(Red, "❌ Error saving profile: %v\n", err)
	}
}

func (y *YuzuChat) loadSystemPrompt() {
	data, err := os.ReadFile(y.systemFile)
	if err != nil {
		if os.IsNotExist(err) {
			colorPrint(Yellow, "📝 No system.txt found, starting without system prompt\n")
			y.systemPrompt = ""
			return
		}
		colorPrint(Red, "❌ Error loading system prompt: %v\n", err)
		return
	}
	y.systemPrompt = string(data)
	lines := strings.Count(y.systemPrompt, "\n") + 1
	chars := len(y.systemPrompt)
	colorPrint(Green, "📖 System prompt loaded from system.txt (%d lines, %d characters)\n", lines, chars)
}

func (y *YuzuChat) saveSystemPrompt(prompt string) error {
	err := os.WriteFile(y.systemFile, []byte(prompt), 0644)
	if err != nil {
		return err
	}
	y.systemPrompt = prompt
	return nil
}

func (y *YuzuChat) loadHistory() {
	data, err := os.ReadFile(y.historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			colorPrint(Yellow, "📝 Starting new conversation history\n")
			y.conversationHistory = []Message{}
			return
		}
		colorPrint(Red, "❌ Error loading history: %v\n", err)
		y.conversationHistory = []Message{}
		return
	}
	var historyData struct {
		Conversations []Message `json:"conversations"`
	}
	if err := json.Unmarshal(data, &historyData); err != nil {
		colorPrint(Red, "❌ Error parsing history: %v\n", err)
		y.conversationHistory = []Message{}
		return
	}
	y.conversationHistory = historyData.Conversations
	colorPrint(Green, "📖 Loaded %d previous messages\n", len(y.conversationHistory))
}

func (y *YuzuChat) saveHistory() {
	historyData := struct {
		Metadata struct {
			LastUpdated      string `json:"last_updated"`
			TotalMessages    int    `json:"total_messages"`
			CurrentModel     string `json:"current_model"`
			CurrentProvider  string `json:"current_provider"`
		} `json:"metadata"`
		Conversations []Message `json:"conversations"`
	}{}
	historyData.Metadata.LastUpdated = time.Now().Format(time.RFC3339)
	historyData.Metadata.TotalMessages = len(y.conversationHistory)
	historyData.Metadata.CurrentModel = y.model
	historyData.Metadata.CurrentProvider = y.currentProvider
	historyData.Conversations = y.conversationHistory
	data, err := json.MarshalIndent(historyData, "", "  ")
	if err != nil {
		colorPrint(Red, "❌ Error marshaling history: %v\n", err)
		return
	}
	if err := os.WriteFile(y.historyFile, data, 0644); err != nil {
		colorPrint(Red, "❌ Error saving history: %v\n", err)
	}
}

func (y *YuzuChat) clearHistory() {
	y.conversationHistory = []Message{}
	err := os.Remove(y.historyFile)
	if err != nil && !os.IsNotExist(err) {
		colorPrint(Red, "❌ Error removing history file: %v\n", err)
		return
	}
	colorPrint(Green, "✅ Conversation history cleared\n")
}

func (y *YuzuChat) addToHistory(role, content string) {
	message := Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now().Format(time.RFC3339),
		Model:     y.model,
		Provider:  y.currentProvider,
	}
	y.conversationHistory = append(y.conversationHistory, message)
	if len(y.conversationHistory) > 20 {
		y.conversationHistory = y.conversationHistory[len(y.conversationHistory)-20:]
	}
	y.saveHistory()
}

func (y *YuzuChat) ListProviders() []string {
	var providers []string
	for name, provider := range y.providers {
		if provider.IsEnabled {
			providers = append(providers, name)
		}
	}
	return providers
}

func (y *YuzuChat) ListModels() []string {
	if provider, exists := y.providers[y.currentProvider]; exists {
		return provider.Models
	}
	return []string{}
}

func (y *YuzuChat) ChangeProvider(providerName string) string {
	if provider, exists := y.providers[providerName]; exists {
		if provider.IsEnabled {
			y.currentProvider = providerName
			if len(provider.Models) > 0 {
				y.model = provider.Models[0]
			}
			y.saveProfile()
			return fmt.Sprintf("✅ Provider changed to: %s", providerName)
		}
		return fmt.Sprintf("❌ Provider '%s' is not enabled (no API key in %s)", providerName, provider.KeyFile)
	}
	return fmt.Sprintf("❌ Provider '%s' not found. Use /providers to see available.", providerName)
}

func (y *YuzuChat) ChangeModel(modelName string) string {
	modelNameLower := strings.ToLower(modelName)
	if provider, exists := y.providers[y.currentProvider]; exists {
		for _, availableModel := range provider.Models {
			if strings.Contains(strings.ToLower(availableModel), modelNameLower) {
				y.model = availableModel
				y.saveProfile()
				return fmt.Sprintf("✅ Model changed to: %s", availableModel)
			}
		}
	}
	return fmt.Sprintf("❌ Model '%s' not found. Use /models to see available.", modelName)
}

func (y *YuzuChat) SetAPIKey(providerName, apiKey string) string {
	if apiKey == "" {
		return "❌ API key cannot be empty"
	}
	provider, exists := y.providers[providerName]
	if !exists {
		return fmt.Sprintf("❌ Provider '%s' not found", providerName)
	}
	if err := y.saveKeyFile(provider.KeyFile, apiKey); err != nil {
		return fmt.Sprintf("❌ Failed to save API key: %v", err)
	}
	provider.APIKey = apiKey
	provider.IsEnabled = true
	return fmt.Sprintf("✅ %s API key saved to %s", providerName, provider.KeyFile)
}

func (y *YuzuChat) RemoveAPIKey(providerName string) string {
	provider, exists := y.providers[providerName]
	if !exists {
		return fmt.Sprintf("❌ Provider '%s' not found", providerName)
	}
	if err := y.removeKeyFile(provider.KeyFile); err != nil {
		return fmt.Sprintf("❌ Failed to remove API key: %v", err)
	}
	provider.APIKey = ""
	provider.IsEnabled = false
	if y.currentProvider == providerName {
		for name, p := range y.providers {
			if p.IsEnabled {
				y.currentProvider = name
				break
			}
		}
	}
	return fmt.Sprintf("✅ %s API key removed from %s", providerName, provider.KeyFile)
}

func (y *YuzuChat) EditSystemPrompt() string {
	colorPrint(Yellow, "📝 Opening system.txt for editing...\n")
	colorPrint(Cyan, "💡 Just edit system.txt directly with your favorite editor!\n")
	colorPrint(Cyan, "💡 The file will be reloaded automatically when you save it.\n")
	colorPrint(Green, "📁 File location: %s\n", y.systemFile)
	return "✅ Edit system.txt directly and use /system reload to refresh"
}

func (y *YuzuChat) ReloadSystemPrompt() string {
	y.loadSystemPrompt()
	lines := strings.Count(y.systemPrompt, "\n") + 1
	chars := len(y.systemPrompt)
	return fmt.Sprintf("✅ System prompt reloaded (%d lines, %d characters)", lines, chars)
}

func (y *YuzuChat) ShowSystemPrompt() string {
	if y.systemPrompt == "" {
		return "No system prompt set (system.txt is empty or doesn't exist)"
	}
	lines := strings.Count(y.systemPrompt, "\n") + 1
	chars := len(y.systemPrompt)
	colorPrint(Cyan, "📋 Current system prompt (%d lines, %d characters):\n", lines, chars)
	colorPrint(Cyan, "────────────────────────────────────────\n")
	fmt.Println(y.systemPrompt)
	colorPrint(Cyan, "────────────────────────────────────────\n")
	return fmt.Sprintf("Displayed %d lines, %d characters", lines, chars)
}

func (y *YuzuChat) SendMessage(message string, stream bool) string {
	provider, exists := y.providers[y.currentProvider]
	if !exists || !provider.IsEnabled {
		return fmt.Sprintf("❌ Provider '%s' is not available", y.currentProvider)
	}
	messages := []map[string]string{}
	if y.systemPrompt != "" {
		messages = append(messages, map[string]string{"role": "system", "content": y.systemPrompt})
	}
	for _, msg := range y.conversationHistory {
		messages = append(messages, map[string]string{"role": msg.Role, "content": msg.Content})
	}
	messages = append(messages, map[string]string{"role": "user", "content": message})
	payload := map[string]interface{}{
		"model":       y.model,
		"messages":    messages,
		"temperature": 0.7,
		"max_tokens":  2048,
		"stream":      stream,
	}
	payloadBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", provider.BaseURL, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return fmt.Sprintf("💥 Request creation failed: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+provider.APIKey)
	req.Header.Set("Content-Type", "application/json")
	if y.currentProvider == "openrouter" {
		req.Header.Set("HTTP-Referer", "https://github.com/icedeyes12/hkkm-yuzu")
		req.Header.Set("X-Title", "Yuzu-Prototype")
	}
	client := &http.Client{Timeout: 60 * time.Second}
	if stream {
		return y.streamResponse(req, message)
	}
	fmt.Printf("🔧 Using: %s/%s...\r", y.currentProvider, y.model)
	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("💥 Request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Sprintf("❌ Error %d: %s", resp.StatusCode, string(body))
	}
	var apiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Sprintf("💥 Response parsing failed: %v", err)
	}
	if len(apiResp.Choices) == 0 {
		return "❌ No response from AI"
	}
	aiResponse := apiResp.Choices[0].Message.Content
	y.addToHistory("user", message)
	y.addToHistory("assistant", aiResponse)
	responseTime := time.Since(startTime).Seconds()
	throughput := float64(apiResp.Usage.CompletionTokens) / responseTime
	stats := fmt.Sprintf("⏱️ %.2fs | 📨 %d→%d tokens | 🚀 %.0f t/s",
		responseTime, apiResp.Usage.PromptTokens, apiResp.Usage.CompletionTokens, throughput)
	fmt.Println(stats)
	return aiResponse
}

func (y *YuzuChat) streamResponse(req *http.Request, userMessage string) string {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("💥 Streaming request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Sprintf("❌ Error %d: %s", resp.StatusCode, string(body))
	}
	colorPrint(Cyan, "🤖: ")
	fullResponse := ""
	startTime := time.Now()
	tokensReceived := 0
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := line[6:]
			if data != "[DONE]" {
				var chunk struct {
					Choices []struct {
						Delta struct {
							Content string `json:"content"`
						} `json:"delta"`
					} `json:"choices"`
				}
				if err := json.Unmarshal([]byte(data), &chunk); err == nil {
					if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
						content := chunk.Choices[0].Delta.Content
						fmt.Print(content)
						fullResponse += content
						tokensReceived += len(content) / 4
					}
				}
			}
		}
	}
	responseTime := time.Since(startTime).Seconds()
	throughput := float64(tokensReceived) / responseTime
	fmt.Printf("\n⏱️ %.2fs | 🚀 ~%.0f t/s (estimated)\n", responseTime, throughput)
	y.addToHistory("user", userMessage)
	y.addToHistory("assistant", fullResponse)
	return fullResponse
}

func (y *YuzuChat) ShowInfo() string {
	totalChars := 0
	for _, msg := range y.conversationHistory {
		totalChars += len(msg.Content)
	}
	enabledProviders := 0
	for _, provider := range y.providers {
		if provider.IsEnabled {
			enabledProviders++
		}
	}
	systemLines := 0
	if y.systemPrompt != "" {
		systemLines = strings.Count(y.systemPrompt, "\n") + 1
	}
	return fmt.Sprintf(`
🍊 Yuzu Prototype - HKMM Project
├── Provider: %s (%s)
├── Model: %s
├── System: %d lines (from system.txt)
├── History: %d exchanges
├── Enabled: %d/%d providers
└── Context: ~%d chars
	`, y.currentProvider, y.providers[y.currentProvider].KeyFile, y.model,
		systemLines, len(y.conversationHistory)/2, enabledProviders, len(y.providers), totalChars)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func main() {
	chat := NewYuzuChat("chat_history.json", "profile.json", "system.txt")
	colorPrint(Purple, `
🍊Yuzu Prototype - HKMM Project♨️
================================
Collab bani ganteng with yuzu herself
https://guthib.com/icedeyes12/hkkm-yuzu

Tips: Type /? for help, /info for status

Special commands:
  /?              - show help
  /info           - show current status  
  /key <provider> <api_key> - set API key
  /exit           - quit
	`)
	scanner := bufio.NewScanner(os.Stdin)
	streaming := false
	for {
		colorPrint(Cyan, "\nYou: ")
		if !scanner.Scan() {
			break
		}
		userInput := strings.TrimSpace(scanner.Text())
		if userInput == "" {
			continue
		}
		if strings.HasPrefix(userInput, "/") {
			parts := strings.Fields(userInput[1:])
			if len(parts) == 0 {
				continue
			}
			command := strings.ToLower(parts[0])
			args := parts[1:]
			switch command {
			case "exit", "quit", "bye":
				colorPrint(Green, "Mata ne~! (Goodbye!)\n")
				return
			case "key":
				if len(args) >= 2 {
					provider := args[0]
					apiKey := strings.Join(args[1:], " ")
					colorPrint(Cyan, "%s\n", chat.SetAPIKey(provider, apiKey))
				} else {
					colorPrint(Yellow, "Usage: /key <provider> <api_key>\n")
					colorPrint(Yellow, "Providers: chutes, openrouter, cerebras\n")
				}
				continue
			case "removekey":
				if len(args) >= 1 {
					provider := args[0]
					colorPrint(Cyan, "%s\n", chat.RemoveAPIKey(provider))
				} else {
					colorPrint(Yellow, "Usage: /removekey <provider>\n")
					colorPrint(Yellow, "Providers: chutes, openrouter, cerebras\n")
				}
				continue
			case "system":
				if len(args) == 0 {
					colorPrint(Cyan, "%s\n", chat.EditSystemPrompt())
				} else if len(args) >= 1 && args[0] == "show" {
					chat.ShowSystemPrompt()
				} else if len(args) >= 1 && args[0] == "reload" {
					colorPrint(Cyan, "%s\n", chat.ReloadSystemPrompt())
				} else {
					newPrompt := strings.Join(args, " ")
					if err := chat.saveSystemPrompt(newPrompt); err != nil {
						colorPrint(Red, "Failed to save system prompt: %v\n", err)
					} else {
						colorPrint(Green, "System prompt updated (%d chars)\n", len(newPrompt))
					}
				}
				continue
			case "help", "?":
				colorPrint(Cyan, `Available Commands:
  /key <provider> <api_key> - Set API key for provider
  /removekey <provider>     - Remove API key for provider
  /system <text>            - Set new system prompt inline
  /system show              - Display current system prompt
  /system reload            - Reload system.txt
  /provider <name>          - Switch provider
  /providers                - List available providers
  /model <name>             - Switch model
  /models                   - List available models
  /clear                    - Clear screen
  /clearhistory             - Clear conversation history
  /info                     - Show current status
  /stream                   - Toggle streaming mode
  /exit, /bye               - Exit
  /help, /?                 - Show this help
`)
				continue
			case "providers":
				providers := chat.ListProviders()
				colorPrint(Cyan, "Available providers:\n")
				for _, p := range providers {
					provider := chat.providers[p]
					if p == chat.currentProvider {
						colorPrint(Yellow, "  - %s <- CURRENT (%s)\n", p, provider.KeyFile)
					} else {
						fmt.Printf("  - %s (%s)\n", p, provider.KeyFile)
					}
				}
				continue
			case "provider":
				if len(args) >= 1 {
					colorPrint(Cyan, "%s\n", chat.ChangeProvider(args[0]))
				} else {
					colorPrint(Yellow, "Usage: /provider <provider_name>\n")
				}
				continue
			case "models":
				models := chat.ListModels()
				colorPrint(Cyan, "Available models for %s:\n", chat.currentProvider)
				for _, m := range models {
					if m == chat.model {
						colorPrint(Yellow, "  - %s <- CURRENT\n", m)
					} else {
						fmt.Printf("  - %s\n", m)
					}
				}
				continue
			case "model":
				if len(args) >= 1 {
					colorPrint(Cyan, "%s\n", chat.ChangeModel(strings.Join(args, " ")))
				} else {
					colorPrint(Yellow, "Usage: /model <model_name>\n")
				}
				continue
			case "info":
				colorPrint(Cyan, "%s\n", chat.ShowInfo())
				continue
			case "stream":
				streaming = !streaming
				colorPrint(Yellow, "Streaming: %s\n", map[bool]string{true: "ON", false: "OFF"}[streaming])
				continue
			case "clear":
				clearScreen()
				colorPrint(Cyan, "Screen cleared\n")
				continue
			case "clearhistory":
				chat.clearHistory()
				continue
			default:
				colorPrint(Red, "Unknown command '/%s'. Type /? for help.\n", command)
				continue
			}
		}
		colorPrint(Yellow, "Thinking with %s/%s...\n", chatResp.Choices[0].Message.Content
	y.addToHistory("user", message)
	y.addToHistory("assistant", aiResponse)
	responseTime := time.Since(startTime).Seconds()
	throughput := float64(apiResp.Usage.CompletionTokens) / responseTime
	stats := fmt.Sprintf("⏱️ %.2fs | 📨 %d→%d tokens | 🚀 %.0f t/s",
		responseTime, apiResp.Usage.PromptTokens, apiResp.Usage.CompletionTokens, throughput)
	fmt.Println(stats)
	return aiResponse
}

func (y *YuzuChat) streamResponse(req *http.Request, userMessage string) string {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("💥 Streaming request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Sprintf("❌ Error %d: %s", resp.StatusCode, string(body))
	}
	colorPrint(Cyan, "🤖: ")
	fullResponse := ""
	startTime := time.Now()
	tokensReceived := 0
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := line[6:]
			if data != "[DONE]" {
				var chunk struct {
					Choices []struct {
						Delta struct {
							Content string `json:"content"`
						} `json:"delta"`
					} `json:"choices"`
				}
				if err := json.Unmarshal([]byte(data), &chunk); err == nil {
					if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
						content := chunk.Choices[0].Delta.Content
						fmt.Print(content)
						fullResponse += content
						tokensReceived += len(content) / 4
					}
				}
			}
		}
	}
	responseTime := time.Since(startTime).Seconds()
	throughput := float64(tokensReceived) / responseTime
	fmt.Printf("\n⏱️ %.2fs | 🚀 ~%.0f t/s (estimated)\n", responseTime, throughput)
	y.addToHistory("user", userMessage)
	y.addToHistory("assistant", fullResponse)
	return fullResponse
}

func (y *YuzuChat) ShowInfo() string {
	totalChars := 0
	for _, msg := range y.conversationHistory {
		totalChars += len(msg.Content)
	}
	enabledProviders := 0
	for _, provider := range y.providers {
		if provider.IsEnabled {
			enabledProviders++
		}
	}
	systemLines := 0
	if y.systemPrompt != "" {
		systemLines = strings.Count(y.systemPrompt, "\n") + 1
	}
	return fmt.Sprintf(`
🍊 Yuzu Prototype - HKMM Project
├── Provider: %s (%s)
├── Model: %s
├── System: %d lines (from system.txt)
├── History: %d exchanges
├── Enabled: %d/%d providers
└── Context: ~%d chars
	`, y.currentProvider, y.providers[y.currentProvider].KeyFile, y.model,
		systemLines, len(y.conversationHistory)/2, enabledProviders, len(y.providers), totalChars)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func main() {
	chat := NewYuzuChat("chat_history.json", "profile.json", "system.txt")
	colorPrint(Purple, `
Yuzu Prototype - HKMM Project
================================
Collab bani ganteng with yuzu herself
https://guthib.com/icedeyes12/hkkm-yuzu

Tips: Type /? for help, /info for status

Special commands:
  /?              - show help
  /info           - show current status  
  /key <provider> <api_key> - set API key
  /exit           - quit
	`)
	scanner := bufio.NewScanner(os.Stdin)
	streaming := false
	for {
		colorPrint(Cyan, "\nYou: ")
		if !scanner.Scan() {
			break
		}
		userInput := strings.TrimSpace(scanner.Text())
		if userInput == "" {
			continue
		}
		if strings.HasPrefix(userInput, "/") {
			parts := strings.Fields(userInput[1:])
			if len(parts) == 0 {
				continue
			}
			command := strings.ToLower(parts[0])
			args := parts[1:]
			switch command {
			case "exit", "quit", "bye":
				colorPrint(Green, "Mata ne~! (Goodbye!)\n")
				return
			case "key":
				if len(args) >= 2 {
					provider := args[0]
					apiKey := strings.Join(args[1:], " ")
					colorPrint(Cyan, "%s\n", chat.SetAPIKey(provider, apiKey))
				} else {
					colorPrint(Yellow, "Usage: /key <provider> <api_key>\n")
					colorPrint(Yellow, "Providers: chutes, openrouter, cerebras\n")
				}
				continue
			case "removekey":
				if len(args) >= 1 {
					provider := args[0]
					colorPrint(Cyan, "%s\n", chat.RemoveAPIKey(provider))
				} else {
					colorPrint(Yellow, "Usage: /removekey <provider>\n")
					colorPrint(Yellow, "Providers: chutes, openrouter, cerebras\n")
				}
				continue
			case "system":
				if len(args) == 0 {
					colorPrint(Cyan, "%s\n", chat.EditSystemPrompt())
				} else if len(args) >= 1 && args[0] == "show" {
					chat.ShowSystemPrompt()
				} else if len(args) >= 1 && args[0] == "reload" {
					colorPrint(Cyan, "%s\n", chat.ReloadSystemPrompt())
				} else {
					newPrompt := strings.Join(args, " ")
					if err := chat.saveSystemPrompt(newPrompt); err != nil {
						colorPrint(Red, "Failed to save system prompt: %v\n", err)
					} else {
						colorPrint(Green, "System prompt updated (%d chars)\n", len(newPrompt))
					}
				}
				continue
			case "help", "?":
				colorPrint(Cyan, `Available Commands:
  /key <provider> <api_key> - Set API key for provider
  /removekey <provider>     - Remove API key for provider
  /system <text...>         - Set new system prompt
  /system show              - Show current system prompt
  /system reload            - Reload system.txt
  /provider <name>          - Switch provider (chutes, openrouter, cerebras)
  /providers                - List available providers
  /model <name>             - Switch model
  /models                   - List available models
  /clear                    - Clear screen
  /clearhistory             - Clear conversation history
  /info                     - Show current status
  /stream                   - Toggle streaming mode
  /exit, /bye               - Exit
  /help, /?                 - Show this help
`)
				continue
			case "providers":
				providers := chat.ListProviders()
				colorPrint(Cyan, "Available providers:\n")
				for _, p := range providers {
					provider := chat.providers[p]
					if p == chat.currentProvider {
						colorPrint(Yellow, "  - %s <- CURRENT (%s)\n", p, provider.KeyFile)
					} else {
						fmt.Printf("  - %s (%s)\n", p, provider.KeyFile)
					}
				}
				continue
			case "provider":
				if len(args) >= 1 {
					colorPrint(Cyan, "%s\n", chat.ChangeProvider(args[0]))
				} else {
					colorPrint(Yellow, "Usage: /provider <provider_name>\n")
				}
				continue
			case "models":
				models := chat.ListModels()
				colorPrint(Cyan, "Available models for %s:\n", chat.currentProvider)
				for _, m := range models {
					if m == chat.model {
						colorPrint(Yellow, "  - %s <- CURRENT\n", m)
					} else {
						fmt.Printf("  - %s\n", m)
					}
				}
				continue
			case "model":
				if len(args) >= 1 {
					colorPrint(Cyan, "%s\n", chat.ChangeModel(strings.Join(args, " ")))
				} else {
					colorPrint(Yellow, "Usage: /model <model_name>\n")
				}
				continue
			case "info":
				colorPrint(Cyan, "%s\n", chat.ShowInfo())
				continue
			case "stream":
				streaming = !streaming
				colorPrint(Yellow, "Streaming: %s\n", map[bool]string{true: "ON", false: "OFF"}[streaming])
				continue
			case "clear":
				clearScreen()
				colorPrint(Cyan, "Screen cleared\n")
				continue
			case "clearhistory":
				chat.clearHistory()
				continue
			default:
				colorPrint(Red, "Unknown command '/%s'. Type /? for help.\n", command)
				continue
			}
		}
		colorPrint(Yellow, "Thinking with %s/%s...\n", chat.currentProvider, chat.model)
		response := chat.SendMessage(userInput, streaming)
		if !streaming {
			colorPrint(Green, "AI: %s\n", response)
		}
	}
}

// titit 𓀐 𓂸