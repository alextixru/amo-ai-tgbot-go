package entities

import (
	"context"
	"fmt"
	"strings"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/filters"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// Service определяет бизнес-логику для работы с основными сущностями amoCRM.
type Service interface {
	// Leads
	SearchLeads(ctx context.Context, filter *gkitmodels.EntitiesFilter, with []string) (*SearchResult, error)
	GetLead(ctx context.Context, id int, with []string) (*EntityResult, error)
	CreateLead(ctx context.Context, data *gkitmodels.EntityData) (*EntityResult, error)
	CreateLeads(ctx context.Context, dataList []gkitmodels.EntityData) ([]*EntityResult, error)
	UpdateLead(ctx context.Context, id int, data *gkitmodels.EntityData) (*EntityResult, error)
	UpdateLeads(ctx context.Context, dataList []gkitmodels.EntityData) ([]*EntityResult, error)
	SyncLead(ctx context.Context, id int, data *gkitmodels.EntityData) (*EntityResult, error)
	LinkLead(ctx context.Context, id int, target *gkitmodels.LinkTarget) (*LinkResult, error)
	UnlinkLead(ctx context.Context, id int, target *gkitmodels.LinkTarget) (*LinkResult, error)

	// Contacts
	SearchContacts(ctx context.Context, filter *gkitmodels.EntitiesFilter, with []string) (*SearchResult, error)
	GetContact(ctx context.Context, id int, with []string) (*EntityResult, error)
	CreateContact(ctx context.Context, data *gkitmodels.EntityData) (*EntityResult, error)
	CreateContacts(ctx context.Context, dataList []gkitmodels.EntityData) ([]*EntityResult, error)
	UpdateContact(ctx context.Context, id int, data *gkitmodels.EntityData) (*EntityResult, error)
	UpdateContacts(ctx context.Context, dataList []gkitmodels.EntityData) ([]*EntityResult, error)
	SyncContact(ctx context.Context, id int, data *gkitmodels.EntityData) (*EntityResult, error)
	LinkContact(ctx context.Context, id int, target *gkitmodels.LinkTarget) (*LinkResult, error)
	UnlinkContact(ctx context.Context, id int, target *gkitmodels.LinkTarget) (*LinkResult, error)
	GetContactChats(ctx context.Context, id int) (any, error)
	LinkContactChats(ctx context.Context, links any) (any, error)

	// Companies
	SearchCompanies(ctx context.Context, filter *gkitmodels.EntitiesFilter, with []string) (*SearchResult, error)
	GetCompany(ctx context.Context, id int, with []string) (*EntityResult, error)
	CreateCompany(ctx context.Context, data *gkitmodels.EntityData) (*EntityResult, error)
	CreateCompanies(ctx context.Context, dataList []gkitmodels.EntityData) ([]*EntityResult, error)
	UpdateCompany(ctx context.Context, id int, data *gkitmodels.EntityData) (*EntityResult, error)
	UpdateCompanies(ctx context.Context, dataList []gkitmodels.EntityData) ([]*EntityResult, error)
	SyncCompany(ctx context.Context, id int, data *gkitmodels.EntityData) (*EntityResult, error)
	LinkCompany(ctx context.Context, id int, target *gkitmodels.LinkTarget) (*LinkResult, error)
	UnlinkCompany(ctx context.Context, id int, target *gkitmodels.LinkTarget) (*LinkResult, error)

	// Метаданные для описания tools
	PipelineNames() []string
	UserNames() []string

	// Справочные данные для Shadow Tool schema response
	StatusesByPipeline() map[string][]string // pipeline_name → []status_name
	LossReasonNames() []string
	CustomFieldCodes(entityType string) []string // "leads"/"contacts"/"companies"
}

type service struct {
	sdk *amocrm.SDK

	// Пользователи
	usersByName map[string]int // name → id
	usersByID   map[int]string // id → name

	// Воронки
	pipelinesByName map[string]int // name → id
	pipelinesByID   map[int]string // id → name

	// Статусы: pipeline_id → status_name → status_id
	statusesByPipelineAndName map[int]map[string]int
	// status_id → status_name (для обратного маппинга)
	statusesByID map[int]string
	// status_id → pipeline_id (для SetStatuses)
	statusPipelineByID map[int]int

	// Кастомные поля по code → field_id (для каждого типа сущности)
	customFieldsLeads    map[string]int // code → id
	customFieldsContacts map[string]int
	customFieldsCompanies map[string]int

	// Причины отказа
	lossReasonsByName map[string]int // name → id
	lossReasonsByID   map[int]string // id → name
}

// New создает новый экземпляр сервиса и загружает справочники.
func New(ctx context.Context, sdk *amocrm.SDK) (Service, error) {
	s := &service{
		sdk:                       sdk,
		usersByName:               make(map[string]int),
		usersByID:                 make(map[int]string),
		pipelinesByName:           make(map[string]int),
		pipelinesByID:             make(map[int]string),
		statusesByPipelineAndName: make(map[int]map[string]int),
		statusesByID:              make(map[int]string),
		statusPipelineByID:        make(map[int]int),
		customFieldsLeads:         make(map[string]int),
		customFieldsContacts:      make(map[string]int),
		customFieldsCompanies:     make(map[string]int),
		lossReasonsByName:         make(map[string]int),
		lossReasonsByID:           make(map[int]string),
	}

	if err := s.loadUsers(ctx); err != nil {
		return nil, fmt.Errorf("entities.New: load users: %w", err)
	}
	if err := s.loadPipelines(ctx); err != nil {
		return nil, fmt.Errorf("entities.New: load pipelines: %w", err)
	}
	if err := s.loadCustomFields(ctx); err != nil {
		return nil, fmt.Errorf("entities.New: load custom fields: %w", err)
	}
	if err := s.loadLossReasons(ctx); err != nil {
		return nil, fmt.Errorf("entities.New: load loss reasons: %w", err)
	}

	return s, nil
}

// loadUsers загружает всех пользователей аккаунта.
func (s *service) loadUsers(ctx context.Context) error {
	usersFilter := filters.NewUsersFilter()
	usersFilter.SetLimit(250)
	users, _, err := s.sdk.Users().Get(ctx, usersFilter)
	if err != nil {
		return err
	}
	for _, u := range users {
		if u.Name != "" {
			s.usersByName[u.Name] = u.ID
			s.usersByID[u.ID] = u.Name
		}
	}
	return nil
}

// loadPipelines загружает воронки и статусы.
func (s *service) loadPipelines(ctx context.Context) error {
	pipelines, _, err := s.sdk.Pipelines().Get(ctx, nil)
	if err != nil {
		return err
	}
	for _, p := range pipelines {
		if p.Name != "" {
			s.pipelinesByName[p.Name] = p.ID
			s.pipelinesByID[p.ID] = p.Name
		}
		if p.Embedded != nil {
			statuses := p.Embedded.Statuses
			if _, ok := s.statusesByPipelineAndName[p.ID]; !ok {
				s.statusesByPipelineAndName[p.ID] = make(map[string]int)
			}
			for _, st := range statuses {
				if st.Name != "" {
					s.statusesByPipelineAndName[p.ID][st.Name] = st.ID
					s.statusesByID[st.ID] = st.Name
					s.statusPipelineByID[st.ID] = p.ID
				}
			}
		}
	}
	return nil
}

// loadCustomFields загружает кастомные поля для leads, contacts, companies.
func (s *service) loadCustomFields(ctx context.Context) error {
	entityTypes := []struct {
		name string
		dest map[string]int
	}{
		{"leads", s.customFieldsLeads},
		{"contacts", s.customFieldsContacts},
		{"companies", s.customFieldsCompanies},
	}

	for _, et := range entityTypes {
		cfFilter := filters.NewCustomFieldsFilter()
		cfFilter.SetLimit(250)
		fields, _, err := s.sdk.CustomFields().Get(ctx, et.name, cfFilter)
		if err != nil {
			// не фатально — продолжаем без кастомных полей для этого типа
			continue
		}
		for _, f := range fields {
			if f.Code != "" {
				et.dest[f.Code] = f.ID
			}
		}
	}
	return nil
}

// loadLossReasons загружает причины отказа.
func (s *service) loadLossReasons(ctx context.Context) error {
	reasons, _, err := s.sdk.LossReasons().Get(ctx, nil) //nolint:staticcheck
	if err != nil {
		// не фатально
		return nil
	}
	for _, r := range reasons {
		if r.Name != "" {
			s.lossReasonsByName[r.Name] = r.ID
			s.lossReasonsByID[r.ID] = r.Name
		}
	}
	return nil
}

// PipelineNames возвращает список названий воронок (для описания tools).
func (s *service) PipelineNames() []string {
	names := make([]string, 0, len(s.pipelinesByName))
	for name := range s.pipelinesByName {
		names = append(names, name)
	}
	return names
}

// UserNames возвращает список имён пользователей (для описания tools).
func (s *service) UserNames() []string {
	names := make([]string, 0, len(s.usersByName))
	for name := range s.usersByName {
		names = append(names, name)
	}
	return names
}

// StatusesByPipeline возвращает маппинг pipeline_name → []status_name.
func (s *service) StatusesByPipeline() map[string][]string {
	result := make(map[string][]string, len(s.pipelinesByID))
	for pipelineID, pipelineName := range s.pipelinesByID {
		statuses, ok := s.statusesByPipelineAndName[pipelineID]
		if !ok {
			result[pipelineName] = []string{}
			continue
		}
		names := make([]string, 0, len(statuses))
		for statusName := range statuses {
			names = append(names, statusName)
		}
		result[pipelineName] = names
	}
	return result
}

// LossReasonNames возвращает список названий причин отказа.
func (s *service) LossReasonNames() []string {
	names := make([]string, 0, len(s.lossReasonsByName))
	for name := range s.lossReasonsByName {
		names = append(names, name)
	}
	return names
}

// CustomFieldCodes возвращает список кодов кастомных полей для указанного типа сущности.
func (s *service) CustomFieldCodes(entityType string) []string {
	var src map[string]int
	switch entityType {
	case "leads":
		src = s.customFieldsLeads
	case "contacts":
		src = s.customFieldsContacts
	case "companies":
		src = s.customFieldsCompanies
	default:
		return nil
	}
	codes := make([]string, 0, len(src))
	for code := range src {
		codes = append(codes, code)
	}
	return codes
}

// --- Резолверы имя → ID ---

func (s *service) resolvePipelineID(name string) (int, error) {
	if name == "" {
		return 0, nil
	}
	id, ok := s.pipelinesByName[name]
	if !ok {
		available := strings.Join(s.PipelineNames(), ", ")
		return 0, fmt.Errorf("воронка %q не найдена. Доступные: %s", name, available)
	}
	return id, nil
}

func (s *service) resolveUserID(name string) (int, error) {
	if name == "" {
		return 0, nil
	}
	id, ok := s.usersByName[name]
	if !ok {
		available := strings.Join(s.UserNames(), ", ")
		return 0, fmt.Errorf("пользователь %q не найден. Доступные: %s", name, available)
	}
	return id, nil
}

func (s *service) resolveUserIDs(names []string) ([]int, error) {
	if len(names) == 0 {
		return nil, nil
	}
	ids := make([]int, 0, len(names))
	for _, name := range names {
		id, err := s.resolveUserID(name)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *service) resolveStatusID(pipelineName, statusName string) (pipelineID, statusID int, err error) {
	if statusName == "" {
		return 0, 0, nil
	}
	pipelineID, err = s.resolvePipelineID(pipelineName)
	if err != nil {
		return 0, 0, err
	}
	statuses, ok := s.statusesByPipelineAndName[pipelineID]
	if !ok {
		return 0, 0, fmt.Errorf("воронка %q не содержит статусов", pipelineName)
	}
	statusID, ok = statuses[statusName]
	if !ok {
		available := make([]string, 0, len(statuses))
		for n := range statuses {
			available = append(available, n)
		}
		return 0, 0, fmt.Errorf("статус %q не найден в воронке %q. Доступные: %s", statusName, pipelineName, strings.Join(available, ", "))
	}
	return pipelineID, statusID, nil
}

func (s *service) resolveLossReasonID(name string) (int, error) {
	if name == "" {
		return 0, nil
	}
	id, ok := s.lossReasonsByName[name]
	if !ok {
		available := make([]string, 0, len(s.lossReasonsByName))
		for n := range s.lossReasonsByName {
			available = append(available, n)
		}
		return 0, fmt.Errorf("причина отказа %q не найдена. Доступные: %s", name, strings.Join(available, ", "))
	}
	return id, nil
}

// --- Резолверы ID → имя ---

func (s *service) lookupUserName(id int) string {
	if id == 0 {
		return ""
	}
	name, ok := s.usersByID[id]
	if !ok {
		return fmt.Sprintf("[unknown:%d]", id)
	}
	return name
}

func (s *service) lookupPipelineName(id int) string {
	if id == 0 {
		return ""
	}
	name, ok := s.pipelinesByID[id]
	if !ok {
		return fmt.Sprintf("[unknown:%d]", id)
	}
	return name
}

func (s *service) lookupStatusName(id int) string {
	if id == 0 {
		return ""
	}
	name, ok := s.statusesByID[id]
	if !ok {
		return fmt.Sprintf("[unknown:%d]", id)
	}
	return name
}

func (s *service) lookupLossReasonName(id int) string {
	if id == 0 {
		return ""
	}
	name, ok := s.lossReasonsByID[id]
	if !ok {
		return fmt.Sprintf("[unknown:%d]", id)
	}
	return name
}
