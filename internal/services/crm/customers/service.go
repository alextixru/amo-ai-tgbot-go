package customers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alextixru/amocrm-sdk-go"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// Service определяет бизнес-логику для работы с покупателями, бонусами, статусами, транзакциями и сегментами.
type Service interface {
	// Customers
	ListCustomers(ctx context.Context, filter *gkitmodels.CustomerFilter, with []string) (*CustomersListOutput, error)
	GetCustomer(ctx context.Context, id int, with []string) (*CustomerOutput, error)
	CreateCustomers(ctx context.Context, data []*gkitmodels.CustomerData) ([]*CustomerOutput, error)
	UpdateCustomers(ctx context.Context, id int, data []*gkitmodels.CustomerData) ([]*CustomerOutput, error)
	DeleteCustomer(ctx context.Context, id int) error
	LinkCustomer(ctx context.Context, customerID int, entityType string, entityID int) error

	// Bonus Points
	GetBonusPoints(ctx context.Context, customerID int) (*BonusPointsInfo, error)
	EarnBonusPoints(ctx context.Context, customerID int, points int) (*BonusPointsResult, error)
	RedeemBonusPoints(ctx context.Context, customerID int, points int) (*BonusPointsResult, error)

	// Statuses
	ListCustomerStatuses(ctx context.Context, page, limit int) ([]StatusOutput, error)
	GetCustomerStatus(ctx context.Context, id int) (*StatusOutput, error)
	CreateCustomerStatuses(ctx context.Context, names []string) ([]StatusOutput, error)
	UpdateCustomerStatus(ctx context.Context, id int, name string) (*StatusOutput, error)
	DeleteCustomerStatus(ctx context.Context, id int) error

	// Transactions
	ListTransactions(ctx context.Context, customerID int, page, limit int) (*TransactionsListOutput, error)
	CreateTransactions(ctx context.Context, customerID int, price int, comment string, accrueBonus bool) (*TransactionsListOutput, error)
	DeleteTransaction(ctx context.Context, customerID int, transactionID int) error

	// Segments
	ListSegments(ctx context.Context, page, limit int) (*SegmentsListOutput, error)
	GetSegment(ctx context.Context, id int) (*SegmentOutput, error)
	CreateSegments(ctx context.Context, names []string) (*SegmentsListOutput, error)
	DeleteSegment(ctx context.Context, id int) error

	// Meta
	UserNames() []string
	StatusNames() []string
}

// BonusPointsInfo текущие бонусные баллы покупателя
type BonusPointsInfo struct {
	// BonusPoints текущий баланс баллов
	BonusPoints int `json:"bonus_points"`
}

// StatusOutput статус покупателя в читаемом формате
type StatusOutput struct {
	ID    int    `json:"id"`
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
	Sort  int    `json:"sort,omitempty"`
}

type service struct {
	sdk            *amocrm.SDK
	usersByName    map[string]int
	usersByID      map[int]string
	statusesByName map[string]int
	statusesByID   map[int]string
}

// New создает новый экземпляр сервиса покупателей и загружает справочники.
func New(ctx context.Context, sdk *amocrm.SDK) (Service, error) {
	s := &service{
		sdk:            sdk,
		usersByName:    make(map[string]int),
		usersByID:      make(map[int]string),
		statusesByName: make(map[string]int),
		statusesByID:   make(map[int]string),
	}
	if err := s.loadUsers(ctx); err != nil {
		return nil, fmt.Errorf("customers: load users: %w", err)
	}
	if err := s.loadStatuses(ctx); err != nil {
		// 422 означает что статусы покупателей недоступны в этом аккаунте (segments mode) — не блокируем старт
		_ = err
	}
	return s, nil
}

// loadUsers загружает всех пользователей из SDK и строит индексы.
func (s *service) loadUsers(ctx context.Context) error {
	users, _, err := s.sdk.Users().Get(ctx, nil)
	if err != nil {
		return err
	}
	for _, u := range users {
		if u == nil || u.Name == "" {
			continue
		}
		s.usersByID[u.ID] = u.Name
		s.usersByName[u.Name] = u.ID
	}
	return nil
}

// loadStatuses загружает статусы покупателей из SDK и строит индексы.
func (s *service) loadStatuses(ctx context.Context) error {
	statuses, _, err := s.sdk.CustomerStatuses().Get(ctx, nil)
	if err != nil {
		return err
	}
	for _, st := range statuses {
		if st.Name == "" {
			continue
		}
		s.statusesByID[st.ID] = st.Name
		s.statusesByName[st.Name] = st.ID
	}
	return nil
}

// resolveUserName переводит имя пользователя в ID.
// Возвращает ошибку с подсказкой если имя не найдено.
func (s *service) resolveUserName(name string) (int, error) {
	if name == "" {
		return 0, nil
	}
	id, ok := s.usersByName[name]
	if !ok {
		available := make([]string, 0, len(s.usersByName))
		for n := range s.usersByName {
			available = append(available, n)
		}
		return 0, fmt.Errorf("пользователь '%s' не найден. Доступные: %s", name, strings.Join(available, ", "))
	}
	return id, nil
}

// resolveUserNames переводит слайс имён в слайс ID.
func (s *service) resolveUserNames(names []string) ([]int, error) {
	ids := make([]int, 0, len(names))
	for _, name := range names {
		id, err := s.resolveUserName(name)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// resolveUserID переводит ID пользователя в имя.
// При неизвестном ID возвращает "[unknown:ID]".
func (s *service) resolveUserID(id int) string {
	if id == 0 {
		return ""
	}
	if name, ok := s.usersByID[id]; ok {
		return name
	}
	return fmt.Sprintf("[unknown:%d]", id)
}

// resolveStatusName переводит имя статуса покупателя в ID.
func (s *service) resolveStatusName(name string) (int, error) {
	if name == "" {
		return 0, nil
	}
	id, ok := s.statusesByName[name]
	if !ok {
		available := make([]string, 0, len(s.statusesByName))
		for n := range s.statusesByName {
			available = append(available, n)
		}
		return 0, fmt.Errorf("статус покупателя '%s' не найден. Доступные: %s", name, strings.Join(available, ", "))
	}
	return id, nil
}

// resolveStatusNames переводит слайс имён статусов в слайс ID.
func (s *service) resolveStatusNames(names []string) ([]int, error) {
	ids := make([]int, 0, len(names))
	for _, name := range names {
		id, err := s.resolveStatusName(name)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// resolveStatusID переводит ID статуса в имя.
func (s *service) resolveStatusID(id int) string {
	if id == 0 {
		return ""
	}
	if name, ok := s.statusesByID[id]; ok {
		return name
	}
	return fmt.Sprintf("[unknown:%d]", id)
}

// parseISO парсит ISO 8601 строку в Unix timestamp.
// Возвращает 0 и ошибку если строка непустая, но не парсится.
func parseISO(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		// Попробуем дату без времени
		t, err = time.Parse("2006-01-02", s)
		if err != nil {
			return 0, fmt.Errorf("неверный формат даты '%s', ожидается ISO 8601 (например 2024-06-01T00:00:00Z или 2024-06-01)", s)
		}
	}
	return t.Unix(), nil
}

// toISO конвертирует Unix timestamp в ISO 8601. Возвращает "" если ts==0.
func toISO(ts int64) string {
	if ts == 0 {
		return ""
	}
	return time.Unix(ts, 0).UTC().Format(time.RFC3339)
}

// UserNames возвращает список доступных имён пользователей.
func (s *service) UserNames() []string {
	names := make([]string, 0, len(s.usersByName))
	for name := range s.usersByName {
		names = append(names, name)
	}
	return names
}

// StatusNames возвращает список доступных имён статусов покупателей.
func (s *service) StatusNames() []string {
	names := make([]string, 0, len(s.statusesByName))
	for name := range s.statusesByName {
		names = append(names, name)
	}
	return names
}
