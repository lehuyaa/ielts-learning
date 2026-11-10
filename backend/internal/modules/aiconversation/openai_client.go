package aiconversation

import "ielts-learning/backend/internal/shared/ai"

// The OpenAI client lives in internal/shared/ai; these aliases keep this
// module's code unchanged.
type (
	openAIMessage = ai.Message
	openAIClient  = ai.Client
)

var ErrAIChatNotConfigured = ai.ErrNotConfigured
