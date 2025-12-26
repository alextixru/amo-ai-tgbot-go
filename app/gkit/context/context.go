// Package context provides user context for AI agents.
package context

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amocrm "github.com/alextixru/amocrm-sdk-go"
)

// DefaultUserID — hardcoded test user ID
const DefaultUserID = 13354958

// UserContext — raw user and account data
type UserContext struct {
	User    map[string]any `json:"user"`
	Account map[string]any `json:"account"`
}

// ToMap returns raw data for AI
func (uc *UserContext) ToMap() map[string]any {
	return map[string]any{
		"user":    uc.User,
		"account": uc.Account,
	}
}

// Builder builds context from amoCRM
type Builder struct {
	sdk *amocrm.SDK
}

// NewBuilder creates builder
func NewBuilder(sdk *amocrm.SDK) *Builder {
	return &Builder{sdk: sdk}
}

// Build gets user context from amoCRM
func (b *Builder) Build(ctx context.Context, telegramUserID int64) (*UserContext, error) {
	amoUserID := DefaultUserID

	user, err := b.sdk.Users().GetOne(ctx, amoUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	account, err := b.sdk.Account().GetCurrent(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Convert to raw JSON
	userJSON, _ := json.Marshal(user)
	accountJSON, _ := json.Marshal(account)

	var userMap, accountMap map[string]any
	json.Unmarshal(userJSON, &userMap)
	json.Unmarshal(accountJSON, &accountMap)

	return &UserContext{User: userMap, Account: accountMap}, nil
}

// MustBuild builds or returns empty context
func (b *Builder) MustBuild(ctx context.Context, telegramUserID int64) *UserContext {
	uc, err := b.Build(ctx, telegramUserID)
	if err != nil {
		log.Printf("[Context] Error: %v", err)
		return &UserContext{User: map[string]any{}, Account: map[string]any{}}
	}
	log.Printf("[Context] Passed")
	return uc
}
