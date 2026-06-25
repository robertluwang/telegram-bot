package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func main() {
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	liteLLMURL := os.Getenv("LITELLM_URL")
	if liteLLMURL == "" {
		log.Fatal("LITELLM_URL is not set (e.g., http://100.x.x.x:4000/v1/chat/completions)")
	}

	liteLLMKey := os.Getenv("LITELLM_KEY")
	if liteLLMKey == "" {
		log.Fatal("LITELLM_KEY is not set")
	}

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() && update.Message.Text != "" {
			go handleMessage(bot, update.Message, liteLLMURL, liteLLMKey)
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, apiURL string, apiKey string) {
	// Let the user know we're typing
	action := tgbotapi.NewChatAction(msg.Chat.ID, tgbotapi.ChatTyping)
	bot.Send(action)

	reqBody := ChatRequest{
		Model: "gemini", // LiteLLM handles the actual routing based on model name, you can change this or pass it via ENV
		Messages: []Message{
			{
				Role:    "user",
				Content: msg.Text,
			},
		},
	}

	jsonValue, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonValue))
	if err != nil {
		sendReply(bot, msg, "Error creating request to AI backend.")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		sendReply(bot, msg, "Error communicating with AI backend.")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("LiteLLM Error: %s", string(bodyBytes))
		sendReply(bot, msg, fmt.Sprintf("AI backend returned status %d", resp.StatusCode))
		return
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		sendReply(bot, msg, "Error decoding AI response.")
		return
	}

	if len(chatResp.Choices) > 0 {
		sendReply(bot, msg, chatResp.Choices[0].Message.Content)
	} else {
		sendReply(bot, msg, "No response from AI.")
	}
}

func sendReply(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, text string) {
	msgReply := tgbotapi.NewMessage(msg.Chat.ID, text)
	msgReply.ReplyToMessageID = msg.MessageID
	if _, err := bot.Send(msgReply); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
