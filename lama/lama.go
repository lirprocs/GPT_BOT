package lama

import (
	"GPT_BOT/config"
	"log"

	"github.com/jpoz/groq"
)

var (
	GROQAPIKey  string
	ModelID     string
	maxMessages = 6
	MaxRequest  = 30
)

var (
	client   *groq.Client
	messages []groq.Message
)

func init() {
	conf := config.New()
	ModelID = conf.ModelID
	GROQAPIKey = conf.GROQAPIKey
	client = groq.NewClient(groq.WithAPIKey(GROQAPIKey))
}

func ProcessMessage(text string) string {
	messages = append(messages, groq.Message{
		Role:    "user",
		Content: text,
	})
	if len(messages) > maxMessages {
		messages = messages[len(messages)-maxMessages:]
	}

	chatCompletion, err := client.CreateChatCompletion(groq.CompletionCreateParams{
		Model:       ModelID,
		Messages:    messages,
		Temperature: 0,
		//Stream:      true,
	})
	if err != nil {
		log.Printf("Error creating chat completion: %v", err)
		return "Error when accessing the server"
	}

	//for delta := range chatCompletion.Stream {
	//	fmt.Print(delta.Choices[0].Delta.Content)
	//}

	response := chatCompletion.Choices[0].Message.Content
	return response
}

func CleanMessage() {
	messages = []groq.Message{}
}

func SetAPIKey(newKey string) {
	GROQAPIKey = newKey
	client = groq.NewClient(groq.WithAPIKey(GROQAPIKey))
}
