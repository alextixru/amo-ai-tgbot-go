package unsorted

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// UnsortedMetadataOutput читаемые метаданные с RFC3339 вместо Unix timestamps.
type UnsortedMetadataOutput struct {
	From           string `json:"from,omitempty"`
	To             string `json:"to,omitempty"`
	Phone          string `json:"phone,omitempty"`
	CalledAt       string `json:"called_at,omitempty"`
	Duration       int    `json:"duration,omitempty"`
	Link           string `json:"link,omitempty"`
	ServiceCode    string `json:"service_code,omitempty"`
	IsCallEvent    bool   `json:"is_call_event,omitempty"`
	Uniq           string `json:"uniq,omitempty"`
	Subject        string `json:"subject,omitempty"`
	ThreadID       string `json:"thread_id,omitempty"`
	MessageID      string `json:"message_id,omitempty"`
	ReceivedAt     string `json:"received_at,omitempty"`
	ContentSummary string `json:"content_summary,omitempty"`
	FormID         string `json:"form_id,omitempty"`
	FormName       string `json:"form_name,omitempty"`
	Page           string `json:"page,omitempty"`
	IP             string `json:"ip,omitempty"`
	Referer        string `json:"referer,omitempty"`
}

// UnsortedOutput читаемое представление одного неразобранного для LLM.
type UnsortedOutput struct {
	UID          string                  `json:"uid,omitempty"`
	Category     string                  `json:"category,omitempty"`
	PipelineName string                  `json:"pipeline_name,omitempty"`
	CreatedAt    string                  `json:"created_at,omitempty"`
	SourceName   string                  `json:"source_name,omitempty"`
	SourceUID    string                  `json:"source_uid,omitempty"`
	Metadata     *UnsortedMetadataOutput `json:"metadata,omitempty"`
	Embedded     *UnsortedEmbeddedOutput `json:"embedded,omitempty"`
}

// UnsortedEmbeddedOutput вложенные данные неразобранного.
type UnsortedEmbeddedOutput struct {
	Leads     []models.Lead    `json:"leads,omitempty"`
	Contacts  []models.Contact `json:"contacts,omitempty"`
	Companies []models.Company `json:"companies,omitempty"`
}

// UnsortedListOutput результат списка неразобранного с пагинацией.
type UnsortedListOutput struct {
	Items    []*UnsortedOutput  `json:"items"`
	PageMeta *services.PageMeta `json:"page_meta,omitempty"`
}

// UnsortedActionResult результат операций accept/decline/link.
type UnsortedActionResult struct {
	UID     string `json:"uid"`
	Success bool   `json:"success"`
}

// Service определяет бизнес-логику для работы с неразобранным.
type Service interface {
	ListUnsorted(ctx context.Context, filter *gkitmodels.UnsortedFilter) (*UnsortedListOutput, error)
	GetUnsorted(ctx context.Context, uid string) (*UnsortedOutput, error)
	CreateUnsorted(ctx context.Context, category string, items []gkitmodels.UnsortedCreateItem) ([]*UnsortedOutput, error)
	AcceptUnsorted(ctx context.Context, uid string, params *gkitmodels.UnsortedAcceptParams) (*UnsortedActionResult, error)
	DeclineUnsorted(ctx context.Context, uid string, params *gkitmodels.UnsortedDeclineParams) (*UnsortedActionResult, error)
	LinkUnsorted(ctx context.Context, uid string, leadID int) (*UnsortedActionResult, error)
	SummaryUnsorted(ctx context.Context, filter *gkitmodels.UnsortedFilter) (*models.UnsortedSummary, error)

	// Meta — для описаний tools
	PipelineNames() []string
	UserNames() []string
	// StatusNames возвращает карту pipeline_name → []status_name для available_values
	StatusNames() map[string][]string
}

type service struct {
	sdk *amocrm.SDK

	usersByName map[string]int
	usersByID   map[int]string

	pipelinesByName map[string]int
	pipelinesByID   map[int]string

	// statusesByPipelineAndName[pipelineID][statusName] = statusID
	statusesByPipelineAndName map[int]map[string]int
}

// New создает новый экземпляр сервиса неразобранного и загружает справочники.
func New(ctx context.Context, sdk *amocrm.SDK) (Service, error) {
	s := &service{
		sdk:                       sdk,
		usersByName:               make(map[string]int),
		usersByID:                 make(map[int]string),
		pipelinesByName:           make(map[string]int),
		pipelinesByID:             make(map[int]string),
		statusesByPipelineAndName: make(map[int]map[string]int),
	}
	if err := s.loadUsers(ctx); err != nil {
		return nil, fmt.Errorf("unsorted: load users: %w", err)
	}
	if err := s.loadPipelines(ctx); err != nil {
		return nil, fmt.Errorf("unsorted: load pipelines: %w", err)
	}
	return s, nil
}

// loadUsers загружает всех пользователей и строит индексы name↔id.
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

// loadPipelines загружает воронки со статусами и строит индексы name↔id.
func (s *service) loadPipelines(ctx context.Context) error {
	params := url.Values{}
	params.Set("with", "statuses")

	pipelines, _, err := s.sdk.Pipelines().Get(ctx, params)
	if err != nil {
		return err
	}
	for _, p := range pipelines {
		if p == nil || p.Name == "" {
			continue
		}
		s.pipelinesByID[p.ID] = p.Name
		s.pipelinesByName[p.Name] = p.ID

		if p.Embedded != nil && len(p.Embedded.Statuses) > 0 {
			statMap := make(map[string]int, len(p.Embedded.Statuses))
			for _, st := range p.Embedded.Statuses {
				if st.Name != "" {
					statMap[st.Name] = st.ID
				}
			}
			s.statusesByPipelineAndName[p.ID] = statMap
		}
	}
	return nil
}

// resolveUserName переводит имя пользователя в ID.
func (s *service) resolveUserName(name string) (int, error) {
	if name == "" {
		return 0, nil
	}
	id, ok := s.usersByName[name]
	if !ok {
		var available []string
		for n := range s.usersByName {
			available = append(available, n)
		}
		return 0, fmt.Errorf("пользователь '%s' не найден. Доступные: %s", name, strings.Join(available, ", "))
	}
	return id, nil
}

// resolveUserID переводит ID пользователя в имя. При неизвестном ID возвращает "[unknown:ID]".
func (s *service) resolveUserID(id int) string {
	if id == 0 {
		return ""
	}
	if name, ok := s.usersByID[id]; ok {
		return name
	}
	return fmt.Sprintf("[unknown:%d]", id)
}

// resolvePipelineName переводит имя воронки в ID.
func (s *service) resolvePipelineName(name string) (int, error) {
	if name == "" {
		return 0, nil
	}
	id, ok := s.pipelinesByName[name]
	if !ok {
		var available []string
		for n := range s.pipelinesByName {
			available = append(available, n)
		}
		return 0, fmt.Errorf("воронка '%s' не найдена. Доступные: %s", name, strings.Join(available, ", "))
	}
	return id, nil
}

// resolvePipelineID переводит ID воронки в имя. При неизвестном ID возвращает "[unknown:ID]".
func (s *service) resolvePipelineID(id int) string {
	if id == 0 {
		return ""
	}
	if name, ok := s.pipelinesByID[id]; ok {
		return name
	}
	return fmt.Sprintf("[unknown:%d]", id)
}

// resolveStatusName переводит имя статуса в ID для заданной воронки.
func (s *service) resolveStatusName(pipelineID int, statusName string) (int, error) {
	if statusName == "" {
		return 0, nil
	}
	statMap, ok := s.statusesByPipelineAndName[pipelineID]
	if !ok {
		return 0, fmt.Errorf("статусы для воронки ID=%d не найдены", pipelineID)
	}
	id, ok := statMap[statusName]
	if !ok {
		var available []string
		for n := range statMap {
			available = append(available, n)
		}
		return 0, fmt.Errorf("статус '%s' не найден. Доступные: %s", statusName, strings.Join(available, ", "))
	}
	return id, nil
}

// PipelineNames возвращает список доступных имён воронок (для описаний tools).
func (s *service) PipelineNames() []string {
	names := make([]string, 0, len(s.pipelinesByName))
	for name := range s.pipelinesByName {
		names = append(names, name)
	}
	return names
}

// UserNames возвращает список доступных имён пользователей (для описаний tools).
func (s *service) UserNames() []string {
	names := make([]string, 0, len(s.usersByName))
	for name := range s.usersByName {
		names = append(names, name)
	}
	return names
}

// StatusNames возвращает карту pipeline_name → []status_name для available_values.
func (s *service) StatusNames() map[string][]string {
	result := make(map[string][]string, len(s.statusesByPipelineAndName))
	for pipelineID, statMap := range s.statusesByPipelineAndName {
		pipelineName := s.resolvePipelineID(pipelineID)
		if pipelineName == "" {
			continue
		}
		statuses := make([]string, 0, len(statMap))
		for name := range statMap {
			statuses = append(statuses, name)
		}
		result[pipelineName] = statuses
	}
	return result
}

// metadataToOutput конвертирует SDK-метаданные в читаемый вывод с RFC3339 строками.
func metadataToOutput(m *models.UnsortedMetadata) *UnsortedMetadataOutput {
	if m == nil {
		return nil
	}
	out := &UnsortedMetadataOutput{
		From:           m.From,
		To:             m.To,
		Phone:          m.Phone,
		Duration:       m.Duration,
		Link:           m.Link,
		ServiceCode:    m.ServiceCode,
		IsCallEvent:    m.IsCallEvent,
		Uniq:           m.Uniq,
		Subject:        m.Subject,
		ThreadID:       m.ThreadID,
		MessageID:      m.MessageID,
		ContentSummary: m.ContentSummary,
		FormID:         m.FormID,
		FormName:       m.FormName,
		Page:           m.Page,
		IP:             m.IP,
		Referer:        m.Referer,
	}
	if m.CalledAt != 0 {
		out.CalledAt = unixToRFC3339(m.CalledAt)
	}
	if m.ReceivedAt != 0 {
		out.ReceivedAt = unixToRFC3339(m.ReceivedAt)
	}
	return out
}

// unsortedToOutput конвертирует SDK-модель в читаемый вывод для LLM.
func (s *service) unsortedToOutput(u *models.Unsorted) *UnsortedOutput {
	if u == nil {
		return nil
	}
	out := &UnsortedOutput{
		UID:          u.UID,
		Category:     string(u.Category),
		PipelineName: s.resolvePipelineID(u.PipelineID),
		SourceName:   u.SourceName,
		SourceUID:    u.SourceUID,
		Metadata:     metadataToOutput(u.Metadata),
	}
	if u.CreatedAt != 0 {
		out.CreatedAt = unixToRFC3339(u.CreatedAt)
	}
	if u.Embedded != nil {
		out.Embedded = &UnsortedEmbeddedOutput{
			Leads:     u.Embedded.Leads,
			Contacts:  u.Embedded.Contacts,
			Companies: u.Embedded.Companies,
		}
	}
	return out
}
