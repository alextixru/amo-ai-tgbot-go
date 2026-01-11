package telegram

import (
	"context"

	"github.com/go-telegram/bot"
)

// Bot wraps Telegram bot instance
type Bot struct {
	b *bot.Bot
}

// New creates a new Telegram bot
func New(token string, opts ...bot.Option) (*Bot, error) {
	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, err
	}
	return &Bot{b: b}, nil
}

// Start starts the bot with long polling
func (tb *Bot) Start(ctx context.Context) {
	tb.b.Start(ctx)
}

// SetDefaultHandler sets the default message handler
func (tb *Bot) SetDefaultHandler(handler bot.HandlerFunc) {
	tb.b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypePrefix, handler)
}

// Underlying returns the underlying bot.Bot for advanced usage
func (tb *Bot) Underlying() *bot.Bot {
	return tb.b
}
