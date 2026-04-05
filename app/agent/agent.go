package agent

import (
	"context"
	"fmt"
	"strings"

	adkagent "google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"

	"github.com/tihn/amo-ai-tgbot-go/app/agent/prompts"
)

const AppName = "amocrm-bot"

// Agent handles AI processing via ADK Runner.
type Agent struct {
	runner         *runner.Runner
	sessionService session.Service
	adkAgent       adkagent.Agent
}

// NewAgent creates a new AI agent backed by ADK Runner.
// Tools are not connected yet — the agent only conducts dialog via LLM.
func NewAgent(ctx context.Context, llmModel model.LLM) (*Agent, error) {
	adkAgent, err := llmagent.New(llmagent.Config{
		Name:        "crm-assistant",
		Model:       llmModel,
		Description: "amoCRM AI assistant",
		Instruction: prompts.BuildSystemPrompt(),
	})
	if err != nil {
		return nil, fmt.Errorf("NewAgent: create llm agent: %w", err)
	}

	sessionService := session.InMemoryService()

	runnr, err := runner.New(runner.Config{
		AppName:           AppName,
		Agent:             adkAgent,
		SessionService:    sessionService,
		AutoCreateSession: true,
	})
	if err != nil {
		return nil, fmt.Errorf("NewAgent: create runner: %w", err)
	}

	return &Agent{
		runner:         runnr,
		sessionService: sessionService,
		adkAgent:       adkAgent,
	}, nil
}

// Process processes a user message through the ADK Runner.
func (a *Agent) Process(ctx context.Context, sessionID, message string) (string, error) {
	userMsg := genai.NewContentFromText(message, genai.RoleUser)

	var result strings.Builder
	for event, err := range a.runner.Run(ctx, sessionID, sessionID, userMsg, adkagent.RunConfig{}) {
		if err != nil {
			return "", fmt.Errorf("agent run: %w", err)
		}
		if event.Content != nil {
			for _, part := range event.Content.Parts {
				if part.Text != "" {
					result.WriteString(part.Text)
				}
			}
		}
	}

	return result.String(), nil
}

// ADKAgent returns the underlying ADK agent (for web launcher).
func (a *Agent) ADKAgent() adkagent.Agent {
	return a.adkAgent
}

// SessionService returns the session service (for web launcher).
func (a *Agent) SessionService() session.Service {
	return a.sessionService
}
