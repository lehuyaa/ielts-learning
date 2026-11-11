package pronunciation

import (
	"log"
	"strings"

	"ielts-learning/backend/internal/shared/ai"
)

// Scorer is the provider chosen to score recordings.
type Scorer struct {
	Provider string
	client   assessor
}

// ChooseScorer picks the provider from PRONUNCIATION_PROVIDER:
// "azure", "openai", or "auto" (default) = Azure when its key is set,
// otherwise OpenAI when OPENAI_API_KEY is set, otherwise disabled.
func ChooseScorer(setting string, azure *AzureClient, openai *ai.Client) Scorer {
	azureScorer := Scorer{Provider: ProviderAzure}
	if azure != nil {
		azureScorer.client = azure
	}
	openaiScorer := Scorer{Provider: ProviderOpenAI, client: NewOpenAIScorer(openai)}

	switch strings.ToLower(strings.TrimSpace(setting)) {
	case "azure":
		return azureScorer
	case "openai":
		return openaiScorer
	case "", "auto":
		if azure != nil && azure.Configured() {
			return azureScorer
		}
		return openaiScorer
	default:
		log.Printf("pronunciation: unknown PRONUNCIATION_PROVIDER %q, using auto", setting)
		return ChooseScorer("auto", azure, openai)
	}
}
