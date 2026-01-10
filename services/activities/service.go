package activities

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

// Service определяет бизнес-логику для работы с активностями amoCRM.
type Service interface {
	// Tasks
	ListTasks(ctx context.Context, parent *gkitmodels.ParentEntity, filter *gkitmodels.ActivityFilter) ([]*models.Task, error)
	GetTask(ctx context.Context, id int) (*models.Task, error)
	CreateTask(ctx context.Context, parent gkitmodels.ParentEntity, data *gkitmodels.ActivityData) (*models.Task, error)
	UpdateTask(ctx context.Context, id int, data *gkitmodels.ActivityData) (*models.Task, error)
	CompleteTask(ctx context.Context, id int, resultText string) (*models.Task, error)

	// Notes
	ListNotes(ctx context.Context, parent gkitmodels.ParentEntity) ([]*models.Note, error)
	GetNote(ctx context.Context, entityType string, id int) (*models.Note, error)
	CreateNote(ctx context.Context, parent gkitmodels.ParentEntity, data *gkitmodels.ActivityData) (*models.Note, error)
	UpdateNote(ctx context.Context, entityType string, id int, data *gkitmodels.ActivityData) (*models.Note, error)

	// Calls
	CreateCall(ctx context.Context, parent gkitmodels.ParentEntity, data *gkitmodels.ActivityData) (*models.Call, error)

	// Events
	ListEvents(ctx context.Context, parent gkitmodels.ParentEntity) ([]*models.Event, error)
	GetEvent(ctx context.Context, id int) (*models.Event, error)

	// Files
	ListFiles(ctx context.Context, parent gkitmodels.ParentEntity) ([]models.FileLink, error)
	LinkFiles(ctx context.Context, parent gkitmodels.ParentEntity, fileUUIDs []string) ([]models.FileLink, error)
	UnlinkFile(ctx context.Context, parent gkitmodels.ParentEntity, fileUUID string) error

	// Links
	ListLinks(ctx context.Context, parent gkitmodels.ParentEntity) ([]*models.EntityLink, error)
	LinkEntity(ctx context.Context, parent gkitmodels.ParentEntity, target *gkitmodels.LinkTarget) ([]*models.EntityLink, error)
	UnlinkEntity(ctx context.Context, parent gkitmodels.ParentEntity, target *gkitmodels.LinkTarget) error

	// Tags
	ListTags(ctx context.Context, entityType string) ([]*models.Tag, error)
	CreateTag(ctx context.Context, entityType string, name string) (*models.Tag, error)
	DeleteTag(ctx context.Context, entityType string, tagID int) error

	// Subscriptions
	ListSubscriptions(ctx context.Context, parent gkitmodels.ParentEntity) ([]models.Subscription, error)
	Subscribe(ctx context.Context, parent gkitmodels.ParentEntity, userIDs []int) ([]models.Subscription, error)
	Unsubscribe(ctx context.Context, parent gkitmodels.ParentEntity, userID int) error

	// Talks
	ListTalks(ctx context.Context, parent gkitmodels.ParentEntity) ([]*models.Talk, error)
	CloseTalk(ctx context.Context, talkID string) error
}

type service struct {
	sdk *amocrm.SDK
}

// New создает новый экземпляр сервиса активностей.
func New(sdk *amocrm.SDK) Service {
	return &service{
		sdk: sdk,
	}
}
