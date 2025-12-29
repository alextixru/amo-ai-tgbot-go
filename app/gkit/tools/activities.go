package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// ActivitiesInput входные параметры для инструмента activities
type ActivitiesInput struct {
	// Parent родительская сущность
	Parent ParentEntity `json:"parent" jsonschema_description:"Родительская сущность: {type, id}"`

	// Layer тип активности: tasks, notes, calls, events, files, links, tags, subscriptions, talks
	Layer string `json:"layer" jsonschema_description:"Тип: tasks, notes, calls, events, files, links, tags, subscriptions, talks"`

	// Action действие
	Action string `json:"action" jsonschema_description:"Действие: list, get, create, update, complete, link, unlink, subscribe, unsubscribe, close"`

	// ID идентификатор элемента (для get, update)
	ID int `json:"id,omitempty" jsonschema_description:"ID элемента (для get, update)"`

	// Data данные для создания/обновления
	Data *ActivityData `json:"data,omitempty" jsonschema_description:"Данные (для create, update)"`

	// UserIDs ID пользователей (для subscribe)
	UserIDs []int `json:"user_ids,omitempty" jsonschema_description:"ID пользователей (для subscribe)"`

	// UserID ID пользователя (для unsubscribe)
	UserID int `json:"user_id,omitempty" jsonschema_description:"ID пользователя (для unsubscribe)"`

	// FileUUIDs UUID файлов (для files.link)
	FileUUIDs []string `json:"file_uuids,omitempty" jsonschema_description:"UUID файлов (для files.link)"`

	// FileUUID UUID файла (для files.unlink)
	FileUUID string `json:"file_uuid,omitempty" jsonschema_description:"UUID файла (для files.unlink)"`

	// TalkID ID чата (для talks.close)
	TalkID string `json:"talk_id,omitempty" jsonschema_description:"ID чата (для talks.close)"`

	// LinkTo цель для связывания (для links)
	LinkTo *LinkTarget `json:"link_to,omitempty" jsonschema_description:"Цель связывания (для links.link/unlink)"`

	// ResultText текст результата (для tasks.complete)
	ResultText string `json:"result_text,omitempty" jsonschema_description:"Текст результата (для tasks.complete)"`
}

// ParentEntity родительская сущность
type ParentEntity struct {
	Type string `json:"type" jsonschema_description:"Тип: leads, contacts, companies"`
	ID   int    `json:"id" jsonschema_description:"ID сущности"`
}

// ActivityData данные для создания/обновления активностей
type ActivityData struct {
	// Task fields
	Text              string `json:"text,omitempty" jsonschema_description:"Текст задачи/примечания"`
	CompleteTillAt    int64  `json:"complete_till_at,omitempty" jsonschema_description:"Срок задачи (unix timestamp)"`
	TaskTypeID        int    `json:"task_type_id,omitempty" jsonschema_description:"ID типа задачи"`
	ResponsibleUserID int    `json:"responsible_user_id,omitempty" jsonschema_description:"ID ответственного"`

	// Note fields
	NoteType string `json:"note_type,omitempty" jsonschema_description:"Тип примечания: common, call_in, call_out, etc."`

	// Call fields
	Direction  string `json:"direction,omitempty" jsonschema_description:"Направление звонка: inbound, outbound"`
	Duration   int    `json:"duration,omitempty" jsonschema_description:"Длительность звонка (секунды)"`
	Source     string `json:"source,omitempty" jsonschema_description:"Источник звонка"`
	Phone      string `json:"phone,omitempty" jsonschema_description:"Номер телефона"`
	CallResult string `json:"call_result,omitempty" jsonschema_description:"Результат звонка"`
	CallStatus int    `json:"call_status,omitempty" jsonschema_description:"Статус звонка: 1-успех, и т.д."`
	UniqID     string `json:"uniq,omitempty" jsonschema_description:"Уникальный ID звонка"`
	Link       string `json:"link,omitempty" jsonschema_description:"Ссылка на запись звонка"`

	// Tag fields
	TagName string `json:"tag_name,omitempty" jsonschema_description:"Название тега"`
	TagID   int    `json:"tag_id,omitempty" jsonschema_description:"ID тега"`
}

// registerActivitiesTool регистрирует инструмент для работы с активностями
func (r *Registry) registerActivitiesTool() {
	r.addTool(genkit.DefineTool[ActivitiesInput, any](
		r.g,
		"activities",
		"Работа с активностями сущностей amoCRM: задачи (tasks), примечания (notes), звонки (calls), "+
			"события (events), файлы (files), связи (links), теги (tags), подписки (subscriptions), чаты (talks). "+
			"Все активности привязаны к parent сущности (lead, contact, company).",
		func(ctx *ai.ToolContext, input ActivitiesInput) (any, error) {
			return r.handleActivities(ctx.Context, input)
		},
	))
}

func (r *Registry) handleActivities(ctx context.Context, input ActivitiesInput) (any, error) {
	// Валидация parent
	if input.Parent.Type == "" || input.Parent.ID == 0 {
		return nil, fmt.Errorf("parent.type and parent.id are required")
	}

	switch input.Layer {
	case "tasks":
		return r.handleTasks(ctx, input)
	case "notes":
		return r.handleNotes(ctx, input)
	case "calls":
		return r.handleCalls(ctx, input)
	case "events":
		return r.handleEvents(ctx, input)
	case "files":
		return r.handleEntityFiles(ctx, input)
	case "links":
		return r.handleLinks(ctx, input)
	case "tags":
		return r.handleTags(ctx, input)
	case "subscriptions":
		return r.handleSubscriptions(ctx, input)
	case "talks":
		return r.handleTalks(ctx, input)
	default:
		return nil, fmt.Errorf("unknown layer: %s", input.Layer)
	}
}

// ============ TASKS ============

func (r *Registry) handleTasks(ctx context.Context, input ActivitiesInput) (any, error) {
	switch input.Action {
	case "list":
		return r.sdk.Tasks().Get(ctx, &services.TasksFilter{
			Limit: 50,
			Page:  1,
		})
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return r.sdk.Tasks().GetOne(ctx, input.ID)
	case "create":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'create'")
		}
		task := models.Task{
			Text:       input.Data.Text,
			EntityID:   input.Parent.ID,
			EntityType: input.Parent.Type,
		}
		if input.Data.CompleteTillAt > 0 {
			task.CompleteTill = &input.Data.CompleteTillAt
		}
		if input.Data.TaskTypeID > 0 {
			task.TaskTypeID = input.Data.TaskTypeID
		}
		if input.Data.ResponsibleUserID > 0 {
			task.ResponsibleUserID = input.Data.ResponsibleUserID
		}
		return r.sdk.Tasks().Create(ctx, []models.Task{task})
	case "update":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'update'")
		}
		task := models.Task{
			BaseModel: models.BaseModel{ID: input.ID},
		}
		if input.Data.Text != "" {
			task.Text = input.Data.Text
		}
		if input.Data.CompleteTillAt > 0 {
			task.CompleteTill = &input.Data.CompleteTillAt
		}
		return r.sdk.Tasks().Update(ctx, []models.Task{task})
	case "complete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'complete'")
		}
		return r.sdk.Tasks().Complete(ctx, input.ID, input.ResultText)
	default:
		return nil, fmt.Errorf("unknown action for tasks: %s", input.Action)
	}
}

// ============ NOTES ============

func (r *Registry) handleNotes(ctx context.Context, input ActivitiesInput) (any, error) {
	switch input.Action {
	case "list":
		return r.sdk.Notes().GetByParent(ctx, input.Parent.Type, input.Parent.ID, &services.NotesFilter{
			Limit: 50,
			Page:  1,
		})
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return r.sdk.Notes().GetOne(ctx, input.Parent.Type, input.ID)
	case "create":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'create'")
		}
		note := models.Note{
			EntityID: input.Parent.ID,
			Params: &models.NoteParams{
				Text: input.Data.Text,
			},
		}
		if input.Data.NoteType != "" {
			note.NoteType = models.NoteType(input.Data.NoteType)
		} else {
			note.NoteType = models.NoteTypeCommon
		}
		return r.sdk.Notes().Create(ctx, input.Parent.Type, []models.Note{note})
	case "update":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'update'")
		}
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'update'")
		}
		note := models.Note{
			BaseModel: models.BaseModel{ID: input.ID},
			Params: &models.NoteParams{
				Text: input.Data.Text,
			},
		}
		return r.sdk.Notes().Update(ctx, input.Parent.Type, []models.Note{note})
	default:
		return nil, fmt.Errorf("unknown action for notes: %s", input.Action)
	}
}

// ============ CALLS ============

func (r *Registry) handleCalls(ctx context.Context, input ActivitiesInput) (any, error) {
	switch input.Action {
	case "create":
		if input.Data == nil {
			return nil, fmt.Errorf("data is required for action 'create'")
		}
		call := models.Call{
			EntityID:   input.Parent.ID,
			EntityType: input.Parent.Type,
			Duration:   input.Data.Duration,
			Source:     input.Data.Source,
			Phone:      input.Data.Phone,
			CallResult: input.Data.CallResult,
			CallStatus: models.CallStatus(input.Data.CallStatus),
		}
		if input.Data.Direction != "" {
			call.Direction = models.CallDirection(input.Data.Direction)
		}
		if input.Data.UniqID != "" {
			call.Uniq = input.Data.UniqID
		}
		if input.Data.Link != "" {
			call.Link = input.Data.Link
		}
		return r.sdk.Calls().AddOne(ctx, &call)
	default:
		return nil, fmt.Errorf("calls only supports 'create' action (write-only API)")
	}
}

// ============ EVENTS ============

func (r *Registry) handleEvents(ctx context.Context, input ActivitiesInput) (any, error) {
	switch input.Action {
	case "list":
		return r.sdk.Events().Get(ctx, &services.EventsFilter{
			Limit:            50,
			Page:             1,
			FilterByEntity:   []string{input.Parent.Type},
			FilterByEntityID: []int{input.Parent.ID},
		})
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required for action 'get'")
		}
		return r.sdk.Events().GetOne(ctx, input.ID)
	default:
		return nil, fmt.Errorf("events only supports 'list' and 'get' actions (read-only API)")
	}
}

// ============ ENTITY FILES ============

func (r *Registry) handleEntityFiles(ctx context.Context, input ActivitiesInput) (any, error) {
	svc := services.NewEntityFilesService(r.sdk.Client(), input.Parent.Type, input.Parent.ID)

	switch input.Action {
	case "list":
		return svc.Get(ctx, 1, 50)
	case "link":
		if len(input.FileUUIDs) == 0 {
			return nil, fmt.Errorf("file_uuids is required for action 'link'")
		}
		return svc.Link(ctx, input.FileUUIDs)
	case "unlink":
		if input.FileUUID == "" {
			return nil, fmt.Errorf("file_uuid is required for action 'unlink'")
		}
		return nil, svc.Unlink(ctx, input.FileUUID)
	default:
		return nil, fmt.Errorf("unknown action for files: %s", input.Action)
	}
}

// ============ LINKS ============

func (r *Registry) handleLinks(ctx context.Context, input ActivitiesInput) (any, error) {
	switch input.Action {
	case "list":
		return r.sdk.Links().Get(ctx, input.Parent.Type, input.Parent.ID, nil)
	case "link":
		if input.LinkTo == nil {
			return nil, fmt.Errorf("link_to is required for action 'link'")
		}
		link := models.EntityLink{
			ToEntityID:   input.LinkTo.ID,
			ToEntityType: input.LinkTo.Type,
		}
		return r.sdk.Links().Link(ctx, input.Parent.Type, input.Parent.ID, []models.EntityLink{link})
	case "unlink":
		if input.LinkTo == nil {
			return nil, fmt.Errorf("link_to is required for action 'unlink'")
		}
		link := models.EntityLink{
			ToEntityID:   input.LinkTo.ID,
			ToEntityType: input.LinkTo.Type,
		}
		return nil, r.sdk.Links().Unlink(ctx, input.Parent.Type, input.Parent.ID, []models.EntityLink{link})
	default:
		return nil, fmt.Errorf("unknown action for links: %s", input.Action)
	}
}

// ============ TAGS ============

func (r *Registry) handleTags(ctx context.Context, input ActivitiesInput) (any, error) {
	switch input.Action {
	case "list":
		return r.sdk.Tags().Get(ctx, input.Parent.Type, &services.TagsFilter{
			Limit: 50,
			Page:  1,
		})
	case "create":
		if input.Data == nil || input.Data.TagName == "" {
			return nil, fmt.Errorf("data.tag_name is required for action 'create'")
		}
		tag := models.Tag{
			Name: input.Data.TagName,
		}
		return r.sdk.Tags().Create(ctx, input.Parent.Type, []models.Tag{tag})
	default:
		return nil, fmt.Errorf("tags supports 'list' and 'create' actions (delete via entity update)")
	}
}

// ============ SUBSCRIPTIONS ============

func (r *Registry) handleSubscriptions(ctx context.Context, input ActivitiesInput) (any, error) {
	svc := services.NewEntitySubscriptionsService(r.sdk.Client(), input.Parent.Type, input.Parent.ID)

	switch input.Action {
	case "list":
		return svc.Get(ctx, 1, 50)
	case "subscribe":
		if len(input.UserIDs) == 0 {
			return nil, fmt.Errorf("user_ids is required for action 'subscribe'")
		}
		return svc.Subscribe(ctx, input.UserIDs)
	case "unsubscribe":
		if input.UserID == 0 {
			return nil, fmt.Errorf("user_id is required for action 'unsubscribe'")
		}
		return nil, svc.Unsubscribe(ctx, input.UserID)
	default:
		return nil, fmt.Errorf("unknown action for subscriptions: %s", input.Action)
	}
}

// ============ TALKS ============

func (r *Registry) handleTalks(ctx context.Context, input ActivitiesInput) (any, error) {
	switch input.Action {
	case "list":
		return r.sdk.Talks().Get(ctx, &services.TalksFilter{
			Limit:      50,
			Page:       1,
			EntityType: input.Parent.Type,
			EntityID:   input.Parent.ID,
		})
	case "close":
		if input.TalkID == "" {
			return nil, fmt.Errorf("talk_id is required for action 'close'")
		}
		return nil, r.sdk.Talks().Close(ctx, input.TalkID)
	default:
		return nil, fmt.Errorf("unknown action for talks: %s", input.Action)
	}
}
