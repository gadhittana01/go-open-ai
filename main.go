package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

type openAI struct {
	client *openai.Client
}

type Message struct {
	Prompt string `json:"prompt"`
}

func (o *openAI) sendMessage(message Message) (string, error) {
	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:            openai.GPT3Dot5Turbo0301,
			MaxTokens:        3000,
			Temperature:      0,
			TopP:             1,
			FrequencyPenalty: 0.5,
			PresencePenalty:  0,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message.Prompt,
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func main() {
	godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	client := openai.NewClient(apiKey)

	openAPI := &openAI{
		client: client,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		message := Message{}

		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := openAPI.sendMessage(message)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return

	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
	log.Println("Server listen at PORT :8000")
}
