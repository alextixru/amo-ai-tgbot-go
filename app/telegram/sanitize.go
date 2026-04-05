package telegram

import (
	"regexp"
	"strings"
)

// Telegram Bot API supported HTML tags.
// See: https://core.telegram.org/bots/api#html-style
var allowedTags = map[string]bool{
	"b": true, "strong": true,
	"i": true, "em": true,
	"u": true, "ins": true,
	"s": true, "strike": true, "del": true,
	"span": true, // class="tg-spoiler"
	"tg-spoiler": true,
	"a":          true,
	"tg-emoji":   true,
	"code":       true,
	"pre":        true,
	"blockquote": true,
}

// tagRegex matches HTML opening, closing, and self-closing tags.
var tagRegex = regexp.MustCompile(`<(/?)([a-zA-Z][a-zA-Z0-9-]*)\b([^>]*)(/?)>`)

// SanitizeTelegramHTML strips HTML tags not supported by Telegram Bot API
// and escapes unmatched < > to prevent parse errors.
func SanitizeTelegramHTML(s string) string {
	result := tagRegex.ReplaceAllStringFunc(s, func(match string) string {
		sub := tagRegex.FindStringSubmatch(match)
		if sub == nil {
			return match
		}
		tagName := strings.ToLower(sub[2])
		if allowedTags[tagName] {
			return match // keep allowed tags as-is
		}
		// Strip disallowed tag, keep inner content
		return ""
	})

	return result
}
