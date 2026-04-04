package telegram

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot/models"
	"github.com/tihn/amo-ai-tgbot-go/app/gkit"
	infraCRM "github.com/tihn/amo-ai-tgbot-go/internal/infrastructure/crm"
	"github.com/tihn/amo-ai-tgbot-go/internal/services/auth"
)

// Service handles Telegram business logic
type Service struct {
	agent     *gkit.Agent
	crmClient *infraCRM.Client
	auth      *auth.Service
}

// NewService creates a new Telegram service
func NewService(agent *gkit.Agent, crmClient *infraCRM.Client, authService *auth.Service) *Service {
	return &Service{
		agent:     agent,
		crmClient: crmClient,
		auth:      authService,
	}
}

// HandleStart returns the start message with connect button
func (s *Service) HandleStart(telegramUserID int64) (string, *models.InlineKeyboardMarkup) {
	isAuth := s.auth.IsAuthenticated(telegramUserID)

	var buttonText, buttonData string
	if isAuth {
		buttonText = "⚙️ Управление Google"
		buttonData = "auth_panel"
	} else {
		buttonText = "🔗 Подключить Google"
		buttonData = "auth_start"
	}

	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{{
			{Text: buttonText, CallbackData: buttonData},
		}},
	}

	message := `👋 Привет! Я amoCRM AI бот.

📋 Доступные команды:
• /status — проверить подключение к amoCRM
• /account — информация об аккаунте
• /pipelines — список воронок и статусов

💬 Или просто напиши мне что-нибудь — я отвечу через AI!`

	return message, keyboard
}

// === Auth Screens ===

// ShowAuthPanel shows the main auth management panel
func (s *Service) ShowAuthPanel(telegramUserID int64) (string, *models.InlineKeyboardMarkup) {
	isAuth := s.auth.IsAuthenticated(telegramUserID)

	if isAuth {
		// Authorized state - show email
		email := s.auth.GetUserEmail(telegramUserID)
		var accountInfo string
		if email != "" {
			accountInfo = fmt.Sprintf("\n\n📧 <b>%s</b>", email)
		}

		message := fmt.Sprintf(`✅ <b>Google аккаунт подключён</b>%s

Ты можешь использовать AI запросы.`, accountInfo)

		keyboard := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{{Text: "🔄 Переподключить", CallbackData: "auth_start"}},
				{{Text: "❌ Отключить", CallbackData: "auth_disconnect"}},
				{{Text: "⬅️ Назад", CallbackData: "back_main"}},
			},
		}
		return message, keyboard
	}

	// Not authorized state
	message := `🔐 <b>Google аккаунт не подключён</b>

Подключи аккаунт для использования AI.`

	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "🔗 Подключить", CallbackData: "auth_start"}},
			{{Text: "⬅️ Назад", CallbackData: "back_main"}},
		},
	}
	return message, keyboard
}

// ShowAuthWaiting shows the waiting for code screen
func (s *Service) ShowAuthWaiting(telegramUserID, chatID int64) (string, *models.InlineKeyboardMarkup) {
	authURL, err := s.auth.StartAuth(telegramUserID, chatID)
	if err != nil {
		return fmt.Sprintf("❌ Ошибка запуска авторизации:\n%v", err), nil
	}

	message := `🔐 <b>Авторизация Google</b>

1️⃣ Нажми кнопку "Открыть ссылку"
2️⃣ Выбери Google аккаунт
3️⃣ Разреши доступ
4️⃣ Скопируй код со страницы
5️⃣ <b>Отправь мне код сообщением</b>

⏱ Код действителен 5 минут.`

	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "🔓 Открыть ссылку", URL: authURL}},
			{{Text: "❌ Отменить", CallbackData: "auth_cancel"}},
		},
	}
	return message, keyboard
}

// ShowAuthSuccess shows the success screen after authorization
func (s *Service) ShowAuthSuccess() (string, *models.InlineKeyboardMarkup) {
	message := `✅ <b>Google аккаунт успешно подключён!</b>

Теперь AI запросы будут выполняться от твоего имени.`

	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "⬅️ В главное меню", CallbackData: "back_main"}},
		},
	}
	return message, keyboard
}

// ShowAuthCanceled shows the canceled screen
func (s *Service) ShowAuthCanceled() (string, *models.InlineKeyboardMarkup) {
	message := `❌ Авторизация отменена.`

	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "🔗 Попробовать снова", CallbackData: "auth_start"}},
			{{Text: "⬅️ Назад", CallbackData: "back_main"}},
		},
	}
	return message, keyboard
}

// ShowAuthDisconnected shows the disconnected screen
func (s *Service) ShowAuthDisconnected() (string, *models.InlineKeyboardMarkup) {
	message := `✅ Google аккаунт отключён.`

	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "🔗 Подключить снова", CallbackData: "auth_start"}},
			{{Text: "⬅️ Назад", CallbackData: "back_main"}},
		},
	}
	return message, keyboard
}

// === Auth Actions ===

// HandleAuthCode processes the authorization code (called when user sends text while waiting)
func (s *Service) HandleAuthCode(ctx context.Context, telegramUserID int64, code string) (string, *models.InlineKeyboardMarkup) {
	if err := s.auth.CompleteAuth(ctx, telegramUserID, code); err != nil {
		message := fmt.Sprintf("❌ <b>Ошибка авторизации</b>\n\n%v", err)
		keyboard := &models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{
				{{Text: "🔄 Попробовать снова", CallbackData: "auth_start"}},
				{{Text: "⬅️ Назад", CallbackData: "back_main"}},
			},
		}
		return message, keyboard
	}
	return s.ShowAuthSuccess()
}

// CancelAuth cancels the pending auth
func (s *Service) CancelAuth(telegramUserID int64) (string, *models.InlineKeyboardMarkup) {
	_ = s.auth.CancelAuth(telegramUserID)
	return s.ShowAuthCanceled()
}

// Disconnect removes the user's tokens
func (s *Service) Disconnect(telegramUserID int64) (string, *models.InlineKeyboardMarkup) {
	_ = s.auth.Logout(telegramUserID)
	return s.ShowAuthDisconnected()
}

// IsWaitingCode returns true if user is waiting to enter auth code
func (s *Service) IsWaitingCode(telegramUserID int64) bool {
	return s.auth.IsWaitingCode(telegramUserID)
}

// === Legacy command handlers (kept for backward compatibility) ===

// HandleConnect starts the Google OAuth flow (legacy)
func (s *Service) HandleConnect(telegramUserID, chatID int64) (string, *models.InlineKeyboardMarkup) {
	return s.ShowAuthWaiting(telegramUserID, chatID)
}

// HandleAuth completes the OAuth flow with the provided code (legacy)
func (s *Service) HandleAuth(ctx context.Context, telegramUserID int64, code string) string {
	msg, _ := s.HandleAuthCode(ctx, telegramUserID, code)
	return msg
}

// HandleMe returns information about the connected account (legacy)
func (s *Service) HandleMe(telegramUserID int64) string {
	if !s.auth.IsAuthenticated(telegramUserID) {
		return "❌ Google аккаунт не подключён.\n\nИспользуй /connect для авторизации."
	}
	return "✅ Google аккаунт подключён.\n\nДля отключения используй /disconnect"
}

// HandleDisconnect removes the user's tokens (legacy)
func (s *Service) HandleDisconnect(telegramUserID int64) string {
	if err := s.auth.Logout(telegramUserID); err != nil {
		return fmt.Sprintf("❌ Ошибка отключения:\n%v", err)
	}
	return "✅ Google аккаунт отключён."
}

// === CRM Handlers ===

// HandleHealthcheck checks CRM connectivity
func (s *Service) HandleHealthcheck(ctx context.Context) string {
	if err := s.crmClient.Healthcheck(ctx); err != nil {
		return fmt.Sprintf("❌ amoCRM недоступен\n\nОшибка: %v", err)
	}
	return "✅ amoCRM доступен!"
}

// HandleAccount returns account information
func (s *Service) HandleAccount(ctx context.Context) string {
	account, err := s.crmClient.SDK().Account().GetCurrent(ctx, nil)
	if err != nil {
		return fmt.Sprintf("❌ Ошибка получения аккаунта\n\n%v", err)
	}
	return fmt.Sprintf(
		"🏢 Аккаунт: %s\n🆔 ID: %d\n🌐 Subdomain: %s",
		account.Name, account.ID, account.Subdomain,
	)
}

// HandlePipelines returns pipelines information
func (s *Service) HandlePipelines(ctx context.Context) string {
	pipelines, _, err := s.crmClient.SDK().Pipelines().Get(ctx, nil)
	if err != nil {
		return fmt.Sprintf("❌ Ошибка получения воронок\n\n%v", err)
	}
	if len(pipelines) == 0 {
		return "📭 Воронок нет"
	}
	var result string
	for _, p := range pipelines {
		result += fmt.Sprintf("📊 %s (ID: %d)\n", p.Name, p.ID)
		statuses, _, err := s.crmClient.SDK().Statuses(p.ID).Get(ctx, nil)
		if err != nil {
			result += fmt.Sprintf("   ⚠️ Ошибка загрузки статусов: %v\n", err)
			continue
		}
		for i, st := range statuses {
			prefix := "├─"
			if i == len(statuses)-1 {
				prefix = "└─"
			}
			result += fmt.Sprintf("   %s %s (ID: %d)\n", prefix, st.Name, st.ID)
		}
		result += "\n"
	}
	return result
}

// ProcessAI processes a message through the AI agent
func (s *Service) ProcessAI(ctx context.Context, telegramUserID int64, chatID int64, text string) (string, error) {
	sessionID := fmt.Sprintf("tg_%d", chatID)
	return s.agent.Process(ctx, sessionID, text)
}

// IsAuthenticated returns true if the user has a valid Google token
func (s *Service) IsAuthenticated(telegramUserID int64) bool {
	return s.auth.IsAuthenticated(telegramUserID)
}
