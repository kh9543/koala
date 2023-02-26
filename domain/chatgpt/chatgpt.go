package chatgpt

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var API_KEY string

type sendQuestionBody struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
	Temperature int    `json:"temperature"`
	MaxTokens   int    `json:"max_tokens"`
}

type sendQuestionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Text          string `json:"text"`
	Index         int    `json:"index"`
	Logprobs      *int   `json:"logprobs"`
	FininshReason string `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func init() {
	API_KEY = os.Getenv("API_KEY")
}

func SendQuestion(msg string) (string, error) {
	client := &http.Client{}

	body := sendQuestionBody{
		Model:       "text-davinci-003",
		Prompt:      msg,
		Temperature: 0,
		MaxTokens:   1024,
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", bytes.NewReader(jsonBytes))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", API_KEY)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	var response sendQuestionResponse
	json.Unmarshal([]byte(res), &response)
	if len(response.Choices) == 0 {
		return "openapi 不知道要怎麼回答你這個問題，請你檢討", nil
	}
	return response.Choices[0].Text, nil
}
