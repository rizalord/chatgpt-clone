package config

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

func NewGemini(
	config *viper.Viper,
	log *logrus.Logger,
) *genai.Client {
	apiKey := config.GetString("GEMINI_API_KEY")

	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is required")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("failed to create Gemini client: %v", err)
	}

	return client
}