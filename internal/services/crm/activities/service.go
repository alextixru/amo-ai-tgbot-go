package activities

import (
	"context"
	"fmt"
	"strings"

	"github.com/alextixru/amocrm-sdk-go"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// Service определяет бизнес-логику для работы с активностями amoCRM.
type Service interface {
	// Tasks
	ListTasks(ctx context.Context, parent *gkitmodels.ParentEntity, filter *gkitmodels.TasksFilter, with []string) (*TasksListOutput, error)
	GetTask(ctx context.Context, id int, with []string) (*TaskOutput, error)
	CreateTask(ctx context.Context, parent gkitmodels.ParentEntity, data *gkitmodels.TaskData) (*TaskOutput, error)
	CreateTasks(ctx context.Context, parent gkitmodels.ParentEntity, data []gkitmodels.TaskData) (*TasksListOutput, error)
	UpdateTask(ctx context.Context, id int, data *gkitmodels.TaskData) (*TaskOutput, error)
	CompleteTask(ctx context.Context, id int, resultText string) (*TaskOutput, error)

	// Notes
	ListNotes(ctx context.Context, parent gkitmodels.ParentEntity, filter *gkitmodels.NotesFilter, with []string) ([]*NoteOutput, error)
	GetNote(ctx context.Context, entityType string, id int) (*NoteOutput, error)
	CreateNote(ctx context.Context, parent gkitmodels.ParentEntity, data *gkitmodels.NoteData) (*NoteOutput, error)
	CreateNotes(ctx context.Context, parent gkitmodels.ParentEntity, data []gkitmodels.NoteData) ([]*NoteOutput, error)
	UpdateNote(ctx context.Context, entityType string, id int, data *gkitmodels.NoteData) (*NoteOutput, error)

	// Calls
	CreateCall(ctx context.Context, parent gkitmodels.ParentEntity, data *gkitmodels.CallData) (*CallOutput, error)

	// Events
	ListEvents(ctx context.Context, parent *gkitmodels.ParentEntity, filter *gkitmodels.EventsFilter) (*EventsListOutput, error)
	GetEvent(ctx context.Context, id int) (*EventOutput, error)

	// Files
	ListFiles(ctx context.Context, parent gkitmodels.ParentEntity, filter *gkitmodels.FilesFilter) (*FilesListOutput, error)
	LinkFiles(ctx context.Context, parent gkitmodels.ParentEntity, fileUUIDs []string) (*FilesListOutput, error)
	UnlinkFile(ctx context.Context, parent gkitmodels.ParentEntity, fileUUID string) error

	// Links
	ListLinks(ctx context.Context, parent gkitmodels.ParentEntity, filter *gkitmodels.LinksFilter) ([]*LinkOutput, error)
	LinkEntity(ctx context.Context, parent gkitmodels.ParentEntity, target *gkitmodels.LinkTarget) ([]*LinkOutput, error)
	LinkEntities(ctx context.Context, parent gkitmodels.ParentEntity, targets []gkitmodels.LinkTarget) ([]*LinkOutput, error)
	UnlinkEntity(ctx context.Context, parent gkitmodels.ParentEntity, target *gkitmodels.LinkTarget) error

	// Tags
	ListTags(ctx context.Context, entityType string, filter *gkitmodels.TagsFilter) ([]*TagOutput, error)
	CreateTag(ctx context.Context, entityType string, name string) (*TagOutput, error)
	CreateTags(ctx context.Context, entityType string, names []string) ([]*TagOutput, error)
	DeleteTag(ctx context.Context, entityType string, tagID int) error
	DeleteTagByName(ctx context.Context, entityType string, tagName string) error

	// Subscriptions
	ListSubscriptions(ctx context.Context, parent gkitmodels.ParentEntity) (*SubscriptionsListOutput, error)
	Subscribe(ctx context.Context, parent gkitmodels.ParentEntity, userNames []string) (*SubscriptionsListOutput, error)
	Unsubscribe(ctx context.Context, parent gkitmodels.ParentEntity, userName string) error

	// Talks
	GetTalk(ctx context.Context, talkID string) (*TalkOutput, error)
	CloseTalk(ctx context.Context, talkID string, forceClose bool) error

	// Meta
	UserNames() []string
}

type service struct {
	sdk         *amocrm.SDK
	usersByName map[string]int
	usersByID   map[int]string
}

// New создает новый экземпляр сервиса активностей и загружает справочники.
func New(ctx context.Context, sdk *amocrm.SDK) (Service, error) {
	s := &service{
		sdk:         sdk,
		usersByName: make(map[string]int),
		usersByID:   make(map[int]string),
	}
	if err := s.loadUsers(ctx); err != nil {
		return nil, fmt.Errorf("activities: load users: %w", err)
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
		// Регистрируем оба варианта: точное имя -> ID.
		// Дубли обрабатываются в resolveUserName.
		s.usersByName[u.Name] = u.ID
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
		// Ищем частичное совпадение для лучшей подсказки
		var available []string
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

// UserNames возвращает список доступных имён пользователей (для описаний tools).
func (s *service) UserNames() []string {
	names := make([]string, 0, len(s.usersByName))
	for name := range s.usersByName {
		names = append(names, name)
	}
	return names
}
