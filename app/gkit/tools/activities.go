package tools

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/tihn/amo-ai-tgbot-go/internal/models"
)

// RegisterActivitiesTool регистрирует инструмент для работы с активностями
func (r *Registry) RegisterActivitiesTool() {
	r.addTool(genkit.DefineTool[models.ActivitiesInput, any](
		r.g,
		"activities",
		"Работа с активностями сущностей amoCRM: задачи (tasks), примечания (notes), звонки (calls), "+
			"события (events), файлы (files), связи (links), теги (tags), подписки (subscriptions), чаты (talks). "+
			"Все активности привязаны к parent сущности (lead, contact, company).",
		func(ctx *ai.ToolContext, input models.ActivitiesInput) (any, error) {
			return r.handleActivities(ctx.Context, input)
		},
	))
}

func (r *Registry) handleActivities(ctx context.Context, input models.ActivitiesInput) (any, error) {
	// Валидация parent: теперь опционально для некоторых действий или при наличии фильтра
	// Для большинства действий (кроме list) parent всё еще обязателен
	if input.Action != "list" && (input.Parent == nil || input.Parent.Type == "" || input.Parent.ID == 0) {
		if input.Action == "create" || input.Action == "update" || input.Action == "complete" || input.Action == "link" || input.Action == "unlink" || input.Action == "subscribe" || input.Action == "unsubscribe" {
			if input.Action == "create" && (input.Parent == nil || input.Parent.ID == 0) {
				return nil, fmt.Errorf("parent.id is required for create action")
			}
		}
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
		return r.handleFiles(ctx, input)
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

func (r *Registry) handleTasks(ctx context.Context, input models.ActivitiesInput) (any, error) {
	switch input.Action {
	case "list":
		return r.activitiesService.ListTasks(ctx, input.Parent, input.Filter, input.With)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required")
		}
		return r.activitiesService.GetTask(ctx, input.ID, input.With)
	case "create":
		if input.Parent == nil {
			return nil, fmt.Errorf("parent is required for create")
		}
		// Batch create
		if len(input.TasksData) > 0 {
			return r.activitiesService.CreateTasks(ctx, *input.Parent, input.TasksData)
		}
		// Single create
		if input.TaskData == nil {
			return nil, fmt.Errorf("task_data or tasks_data is required")
		}
		return r.activitiesService.CreateTask(ctx, *input.Parent, input.TaskData)
	case "update":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required")
		}
		if input.TaskData == nil {
			return nil, fmt.Errorf("task_data is required")
		}
		return r.activitiesService.UpdateTask(ctx, input.ID, input.TaskData)
	case "complete":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required")
		}
		return r.activitiesService.CompleteTask(ctx, input.ID, input.ResultText)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

func (r *Registry) handleNotes(ctx context.Context, input models.ActivitiesInput) (any, error) {
	if input.Parent == nil {
		return nil, fmt.Errorf("parent is required for notes")
	}
	switch input.Action {
	case "list":
		return r.activitiesService.ListNotes(ctx, *input.Parent, input.NotesFilter, input.With)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required")
		}
		return r.activitiesService.GetNote(ctx, input.Parent.Type, input.ID)
	case "create":
		// Batch create
		if len(input.NotesData) > 0 {
			return r.activitiesService.CreateNotes(ctx, *input.Parent, input.NotesData)
		}
		// Single create
		if input.NoteData == nil {
			return nil, fmt.Errorf("note_data or notes_data is required")
		}
		return r.activitiesService.CreateNote(ctx, *input.Parent, input.NoteData)
	case "update":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required")
		}
		if input.NoteData == nil {
			return nil, fmt.Errorf("note_data is required")
		}
		return r.activitiesService.UpdateNote(ctx, input.Parent.Type, input.ID, input.NoteData)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

func (r *Registry) handleCalls(ctx context.Context, input models.ActivitiesInput) (any, error) {
	if input.Action != "create" {
		return nil, fmt.Errorf("calls only supports 'create' action")
	}
	if input.Parent == nil {
		return nil, fmt.Errorf("parent is required for calls")
	}
	if input.CallData == nil {
		return nil, fmt.Errorf("call_data is required")
	}
	return r.activitiesService.CreateCall(ctx, *input.Parent, input.CallData)
}

func (r *Registry) handleEvents(ctx context.Context, input models.ActivitiesInput) (any, error) {
	switch input.Action {
	case "list":
		return r.activitiesService.ListEvents(ctx, input.Parent, input.EventsFilter)
	case "get":
		if input.ID == 0 {
			return nil, fmt.Errorf("id is required")
		}
		return r.activitiesService.GetEvent(ctx, input.ID)
	default:
		return nil, fmt.Errorf("events only supports 'list' and 'get' actions")
	}
}

func (r *Registry) handleFiles(ctx context.Context, input models.ActivitiesInput) (any, error) {
	if input.Parent == nil {
		return nil, fmt.Errorf("parent is required for files")
	}
	switch input.Action {
	case "list":
		return r.activitiesService.ListFiles(ctx, *input.Parent, input.FilesFilter)
	case "link":
		if len(input.FileUUIDs) == 0 {
			return nil, fmt.Errorf("file_uuids is required")
		}
		return r.activitiesService.LinkFiles(ctx, *input.Parent, input.FileUUIDs)
	case "unlink":
		if input.FileUUID == "" {
			return nil, fmt.Errorf("file_uuid is required")
		}
		return nil, r.activitiesService.UnlinkFile(ctx, *input.Parent, input.FileUUID)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

func (r *Registry) handleLinks(ctx context.Context, input models.ActivitiesInput) (any, error) {
	if input.Parent == nil {
		return nil, fmt.Errorf("parent is required for links")
	}
	switch input.Action {
	case "list":
		return r.activitiesService.ListLinks(ctx, *input.Parent, input.LinksFilter)
	case "link":
		// Batch link
		if len(input.LinksTo) > 0 {
			return r.activitiesService.LinkEntities(ctx, *input.Parent, input.LinksTo)
		}
		// Single link
		if input.LinkTo == nil {
			return nil, fmt.Errorf("link_to or links_to is required")
		}
		return r.activitiesService.LinkEntity(ctx, *input.Parent, input.LinkTo)
	case "unlink":
		if input.LinkTo == nil {
			return nil, fmt.Errorf("link_to is required")
		}
		return nil, r.activitiesService.UnlinkEntity(ctx, *input.Parent, input.LinkTo)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

func (r *Registry) handleTags(ctx context.Context, input models.ActivitiesInput) (any, error) {
	entityType := ""
	if input.Parent != nil {
		entityType = input.Parent.Type
	}
	if entityType == "" {
		return nil, fmt.Errorf("parent.type is required for tags")
	}

	switch input.Action {
	case "list":
		return r.activitiesService.ListTags(ctx, entityType, input.TagsFilter)
	case "create":
		// Batch create
		if len(input.TagNames) > 0 {
			return r.activitiesService.CreateTags(ctx, entityType, input.TagNames)
		}
		// Single create
		if input.TagName == "" {
			return nil, fmt.Errorf("tag_name or tag_names is required")
		}
		return r.activitiesService.CreateTag(ctx, entityType, input.TagName)
	case "delete":
		if input.TagID == 0 {
			return nil, fmt.Errorf("tag_id is required")
		}
		return nil, r.activitiesService.DeleteTag(ctx, entityType, input.TagID)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

func (r *Registry) handleSubscriptions(ctx context.Context, input models.ActivitiesInput) (any, error) {
	if input.Parent == nil {
		return nil, fmt.Errorf("parent is required for subscriptions")
	}
	switch input.Action {
	case "list":
		return r.activitiesService.ListSubscriptions(ctx, *input.Parent)
	case "subscribe":
		if len(input.UserIDs) == 0 {
			return nil, fmt.Errorf("user_ids is required")
		}
		return r.activitiesService.Subscribe(ctx, *input.Parent, input.UserIDs)
	case "unsubscribe":
		if input.UserID == 0 {
			return nil, fmt.Errorf("user_id is required")
		}
		return nil, r.activitiesService.Unsubscribe(ctx, *input.Parent, input.UserID)
	default:
		return nil, fmt.Errorf("unknown action: %s", input.Action)
	}
}

func (r *Registry) handleTalks(ctx context.Context, input models.ActivitiesInput) (any, error) {
	switch input.Action {
	case "close":
		if input.TalkID == "" {
			return nil, fmt.Errorf("talk_id is required")
		}
		return nil, r.activitiesService.CloseTalk(ctx, input.TalkID, input.ForceClose)
	default:
		return nil, fmt.Errorf("unknown action: %s (talks only support 'close')", input.Action)
	}
}
