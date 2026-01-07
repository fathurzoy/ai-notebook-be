package chatbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type GeminiChatParts struct {
	Text string `json:"text"`
}

type GeminiChatContent struct {
	Parts []GeminiChatParts `json:"parts"`
	Role  string            `json:"role"`
}

type GeminiChatRequest struct {
	Contents []GeminiChatContent `json:"contents"`
}

type ChatHistory struct {
	Chat string
	Role string
}

type GeminiChatCandidate struct {
	Content *GeminiChatContent `json:"content"`
}

type GeminiChatResponse struct {
	Candidates []GeminiChatCandidate `json:"candidates"`
}

func GetGeminiResponse(ctx context.Context, apiKey string, chatHistories []*ChatHistory) (string, error) {

	chatContents := make([]GeminiChatContent, 0)
	for i, chatHistory := range chatHistories {
		chatContents[i] = GeminiChatContent{
			Role:  chatHistory.Role,
			Parts: []GeminiChatParts{{Text: chatHistory.Chat}},
		}
	}

	payload := GeminiChatRequest{
		Contents: chatContents,
	}
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash-exp:generateContent",
		bytes.NewBuffer(payloadJson),
	)

	if err != nil {
		return "", err
	}

	req.Header.Set("x-goog-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status error, got status %d. with response body %s", res.StatusCode, string(resBody))
	}

	var geminiRes GeminiChatResponse

	err = json.Unmarshal(resBody, &geminiRes)
	if err != nil {
		return "", err
	}

	log.Println(geminiRes.Candidates[0].Content.Parts[0].Text)

	return geminiRes.Candidates[0].Content.Parts[0].Text, nil
}
