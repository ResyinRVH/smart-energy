package service

import (
	"a21hc3NpZ25tZW50/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}

type HuggingFaceResponse struct {
	Embedding []float64 `json:"embedding"`
}

func (s *AIService) AnalyzeData(table map[string][]string, query, token string) (string, error) {
	// return nil, nil // TODO: replace this
	if len(table) == 0 {
		return "", errors.New("table cannot be empty")
	}

	requestBody := model.AIRequest{
		Inputs: model.Inputs{
			Table: table,
			Query: query,
		},
	}

	payloadBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/google/tapas-base-finetuned-wtq", bytes.NewReader(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error response from AI model: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response model.TapasResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if len(response.Cells) == 0 {
		return "", errors.New("no answer found in response")
	}

	return response.Cells[0], nil
}

func (s *AIService) RecomendationFromLocation(location, energyUse, token string) (string, error) {
	inputText := fmt.Sprintf("For %s location, with energy usage: %s, recomendation is <mask>.", location, energyUse)

	requestBody := map[string]interface{}{
		"inputs": inputText,
	}

	payloadBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to encode payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/FacebookAI/xlm-roberta-base", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to read response: %w", err)
	// }

	var responseBody []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if len(responseBody) == 0 || responseBody[0]["token_str"] == nil {
		return "", fmt.Errorf("unexpected API response format")
	}

	recommendation := responseBody[0]["token_str"].(string)
	return recommendation, nil
}

func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {

	// Membentuk body request berdasarkan format curl
	requestBody := map[string]interface{}{
		"model": "Qwen/Qwen2.5-1.5B-Instruct",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": fmt.Sprintf("%s %s", context, query),
			},
		},
		"max_tokens": 500,
		"stream":     false,
	}

	payloadBytes, err := json.Marshal(requestBody)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Membuat request HTTP
	req, err := http.NewRequest(
		"POST",
		"https://api-inference.huggingface.co/models/microsoft/Phi-3.5-mini-instruct/v1/chat/completions",
		bytes.NewReader(payloadBytes),
	)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to create request: %v", err)
	}

	// Menambahkan header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	// Mengirim request
	resp, err := s.Client.Do(req)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	// Memproses respons
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to decode response: %v", err)
	}

	// Pastikan ada pilihan dan konten
	if len(response.Choices) == 0 || response.Choices[0].Message.Content == "" {
		return model.ChatResponse{}, errors.New("no answer found in response")
	}

	// Mengembalikan hasil
	return model.ChatResponse{
		GeneratedText: response.Choices[0].Message.Content,
	}, nil
}
