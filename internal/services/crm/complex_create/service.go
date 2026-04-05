package complex_create

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/alextixru/amocrm-sdk-go"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// Service определяет бизнес-логику для комплексного создания сущностей (сделка + контакты/компания).
type Service interface {
	// CreateComplex создаёт сделку вместе с контактами и/или компанией.
	// Принимает входные данные с именами (pipeline_name, status_name, responsible_user_name).
	// Возвращает обогащённый ответ с именами вместо числовых ID.
	CreateComplex(ctx context.Context, input *gkitmodels.ComplexCreateInput) (*ComplexCreateResult, error)

	// CreateComplexBatch создаёт несколько сделок за один запрос (до 50 штук).
	CreateComplexBatch(ctx context.Context, inputs []gkitmodels.ComplexCreateInput) ([]ComplexCreateResult, error)

	// PipelineNames возвращает список доступных воронок для использования в описаниях tools.
	PipelineNames() []string

	// UserNames возвращает список доступных пользователей для использования в описаниях tools.
	UserNames() []string

	// StatusesByPipeline возвращает карту pipeline_name → []status_name для schema response.
	StatusesByPipeline() map[string][]string
}

type service struct {
	sdk *amocrm.SDK

	// usersByName имя пользователя → ID
	usersByName map[string]int
	// usersByID ID пользователя → имя
	usersByID map[int]string

	// pipelinesByName имя воронки → ID
	pipelinesByName map[string]int
	// pipelinesByID ID воронки → имя
	pipelinesByID map[int]string

	// statusesByPipelineAndName pipelineID → statusName → statusID
	statusesByPipelineAndName map[int]map[string]int
	// statusesByPipelineAndID pipelineID → statusID → statusName
	statusesByPipelineAndID map[int]map[int]string
}

// New создает новый экземпляр сервиса комплексного создания.
// При инициализации загружает справочники пользователей и воронок из SDK.
// Ошибки загрузки не блокируют старт: сервис будет работать с пустыми справочниками
// и будет возвращать понятные ошибки при попытке резолвинга имён.
func New(ctx context.Context, sdk *amocrm.SDK) (Service, error) {
	svc := &service{
		sdk:                       sdk,
		usersByName:               make(map[string]int),
		usersByID:                 make(map[int]string),
		pipelinesByName:           make(map[string]int),
		pipelinesByID:             make(map[int]string),
		statusesByPipelineAndName: make(map[int]map[string]int),
		statusesByPipelineAndID:   make(map[int]map[int]string),
	}

	if err := svc.loadUsers(ctx); err != nil {
		return nil, fmt.Errorf("complex_create: не удалось загрузить пользователей: %w", err)
	}

	if err := svc.loadPipelines(ctx); err != nil {
		return nil, fmt.Errorf("complex_create: не удалось загрузить воронки: %w", err)
	}

	return svc, nil
}

// loadUsers загружает всех пользователей аккаунта и строит двунаправленные мапы.
func (s *service) loadUsers(ctx context.Context) error {
	users, _, err := s.sdk.Users().Get(ctx, nil)
	if err != nil {
		return err
	}

	for _, u := range users {
		if u == nil || u.Name == "" {
			continue
		}
		s.usersByName[u.Name] = u.ID
		s.usersByID[u.ID] = u.Name
	}

	return nil
}

// loadPipelines загружает все воронки со статусами и строит индексы.
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
		s.pipelinesByName[p.Name] = p.ID
		s.pipelinesByID[p.ID] = p.Name

		if p.Embedded == nil {
			continue
		}

		byName := make(map[string]int)
		byID := make(map[int]string)
		for _, st := range p.Embedded.Statuses {
			if st.Name == "" {
				continue
			}
			byName[st.Name] = st.ID
			byID[st.ID] = st.Name
		}
		s.statusesByPipelineAndName[p.ID] = byName
		s.statusesByPipelineAndID[p.ID] = byID
	}

	return nil
}

// PipelineNames возвращает список имён всех загруженных воронок.
func (s *service) PipelineNames() []string {
	names := make([]string, 0, len(s.pipelinesByName))
	for name := range s.pipelinesByName {
		names = append(names, name)
	}
	return names
}

// UserNames возвращает список имён всех загруженных пользователей.
func (s *service) UserNames() []string {
	names := make([]string, 0, len(s.usersByName))
	for name := range s.usersByName {
		names = append(names, name)
	}
	return names
}

// StatusesByPipeline возвращает карту pipeline_name → []status_name для schema response.
func (s *service) StatusesByPipeline() map[string][]string {
	result := make(map[string][]string, len(s.pipelinesByName))
	for pName, pID := range s.pipelinesByName {
		byName, ok := s.statusesByPipelineAndName[pID]
		if !ok {
			result[pName] = []string{}
			continue
		}
		statuses := make([]string, 0, len(byName))
		for sName := range byName {
			statuses = append(statuses, sName)
		}
		result[pName] = statuses
	}
	return result
}

// resolveUserID возвращает ID пользователя по имени.
// Если имя пустое — возвращает 0, nil (не обязательное поле).
func (s *service) resolveUserID(name string) (int, error) {
	if name == "" {
		return 0, nil
	}
	id, ok := s.usersByName[name]
	if !ok {
		available := make([]string, 0, len(s.usersByName))
		for n := range s.usersByName {
			available = append(available, n)
		}
		return 0, fmt.Errorf("пользователь %q не найден. Доступные: %s", name, strings.Join(available, ", "))
	}
	return id, nil
}

// resolvePipelineID возвращает ID воронки по имени.
// Если имя пустое — возвращает 0, nil (не обязательное поле).
func (s *service) resolvePipelineID(name string) (int, error) {
	if name == "" {
		return 0, nil
	}
	id, ok := s.pipelinesByName[name]
	if !ok {
		available := make([]string, 0, len(s.pipelinesByName))
		for n := range s.pipelinesByName {
			available = append(available, n)
		}
		return 0, fmt.Errorf("воронка %q не найдена. Доступные: %s", name, strings.Join(available, ", "))
	}
	return id, nil
}

// resolveStatusID возвращает ID статуса по имени внутри заданной воронки.
// Если имя пустое — возвращает 0, nil (не обязательное поле).
func (s *service) resolveStatusID(pipelineName string, pipelineID int, statusName string) (int, error) {
	if statusName == "" {
		return 0, nil
	}
	byName, ok := s.statusesByPipelineAndName[pipelineID]
	if !ok || len(byName) == 0 {
		return 0, fmt.Errorf("статус %q не найден: воронка %q (id=%d) не имеет загруженных статусов", statusName, pipelineName, pipelineID)
	}
	id, ok := byName[statusName]
	if !ok {
		available := make([]string, 0, len(byName))
		for n := range byName {
			available = append(available, n)
		}
		return 0, fmt.Errorf("статус %q не найден в воронке %q. Доступные: %s", statusName, pipelineName, strings.Join(available, ", "))
	}
	return id, nil
}

// lookupUserName возвращает имя пользователя по ID.
// Если ID не найден — возвращает "[unknown:ID]".
func (s *service) lookupUserName(id int) string {
	if id == 0 {
		return ""
	}
	if name, ok := s.usersByID[id]; ok {
		return name
	}
	return fmt.Sprintf("[unknown:%d]", id)
}

// lookupPipelineName возвращает имя воронки по ID.
func (s *service) lookupPipelineName(id int) string {
	if id == 0 {
		return ""
	}
	if name, ok := s.pipelinesByID[id]; ok {
		return name
	}
	return fmt.Sprintf("[unknown:%d]", id)
}

// lookupStatusName возвращает имя статуса по ID воронки и ID статуса.
func (s *service) lookupStatusName(pipelineID, statusID int) string {
	if statusID == 0 {
		return ""
	}
	if byID, ok := s.statusesByPipelineAndID[pipelineID]; ok {
		if name, ok := byID[statusID]; ok {
			return name
		}
	}
	return fmt.Sprintf("[unknown:%d]", statusID)
}
