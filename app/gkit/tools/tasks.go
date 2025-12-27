package tools

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/models"
	"github.com/alextixru/amocrm-sdk-go/core/services"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

// TaskToolInput входные параметры для инструмента crm_tasks
type TaskToolInput struct {
	// Action действие: list, create, complete
	Action string `json:"action"`

	// Filter фильтр для поиска (используется при action=list)
	Filter *TaskFilterInput `json:"filter,omitempty"`

	// ID идентификатор задачи (используется при action=complete)
	ID int `json:"id,omitempty"`

	// Data данные для создания или завершения (используется при action=create, complete)
	Data *TaskDataInput `json:"data,omitempty"`
}

// TaskFilterInput параметры фильтрации задач
type TaskFilterInput struct {
	// ResponsibleUserID ID ответственного
	ResponsibleUserID int `json:"responsible_user_id,omitempty"`
	// TaskTypeID ID типа задачи
	TaskTypeID int `json:"task_type_id,omitempty"`
	// IsCompleted статус завершенности (true/false)
	IsCompleted *bool `json:"is_completed,omitempty"`
	// Limit количество элементов (по умолчанию 50)
	Limit int `json:"limit,omitempty"`
}

// TaskDataInput данные задачи
type TaskDataInput struct {
	// Text текст задачи (обязательно при создании)
	Text string `json:"text,omitempty"`
	// CompleteTill дедлайн (timestamp)
	CompleteTill int64 `json:"complete_till,omitempty"`
	// ResponsibleUserID ID ответственного
	ResponsibleUserID int `json:"responsible_user_id,omitempty"`
	// EntityType тип привязки (leads, contacts, companies, customers)
	EntityType string `json:"entity_type,omitempty"`
	// EntityID ID сущности
	EntityID int `json:"entity_id,omitempty"`
	// TaskTypeID ID типа задачи (Call, Meeting и т.д.)
	TaskTypeID int `json:"task_type_id,omitempty"`
	// ResultText текст результата (для complete)
	ResultText string `json:"result_text,omitempty"`
}

// registerTasksTools регистрирует инструменты для работы с задачами
func (r *Registry) registerTasksTools() {
	r.addTool(genkit.DefineTool[TaskToolInput, any](
		r.g,
		"crm_tasks",
		"Управление задачами (Tasks). Позволяет создавать, просматривать и завершать задачи.",
		func(ctx *ai.ToolContext, input TaskToolInput) (any, error) {
			switch input.Action {
			case "list":
				return r.listTasks(ctx.Context, input.Filter)
			case "create":
				if input.Data == nil {
					return nil, fmt.Errorf("data is required for 'create' action")
				}
				return r.createTask(ctx.Context, input.Data)
			case "complete":
				if input.ID == 0 {
					return nil, fmt.Errorf("id is required for 'complete' action")
				}
				resultText := ""
				if input.Data != nil {
					resultText = input.Data.ResultText
				}
				return r.sdk.Tasks().Complete(ctx.Context, input.ID, resultText)
			default:
				return nil, fmt.Errorf("unknown action: %s", input.Action)
			}
		},
	))
}

// listTasks выполняет поиск задач
func (r *Registry) listTasks(ctx context.Context, filter *TaskFilterInput) ([]models.Task, error) {
	sdkFilter := &services.TasksFilter{
		Limit: 50,
	}

	if filter != nil {
		if filter.Limit > 0 {
			sdkFilter.Limit = filter.Limit
		}
		if filter.ResponsibleUserID > 0 {
			sdkFilter.FilterByResponsibleUserID = []int{filter.ResponsibleUserID}
		}
		if filter.TaskTypeID > 0 {
			sdkFilter.FilterByTaskTypeID = []int{filter.TaskTypeID}
		}
		if filter.IsCompleted != nil {
			sdkFilter.FilterByIsCompleted = filter.IsCompleted
		}
	}

	return r.sdk.Tasks().Get(ctx, sdkFilter)
}

// createTask создает новую задачу
func (r *Registry) createTask(ctx context.Context, data *TaskDataInput) (*models.Task, error) {
	task := models.Task{
		Text:         data.Text,
		CompleteTill: &data.CompleteTill,
		TaskTypeID:   data.TaskTypeID,
		EntityType:   data.EntityType,
		EntityID:     data.EntityID,
	}

	if data.ResponsibleUserID > 0 {
		task.ResponsibleUserID = data.ResponsibleUserID
	}

	createdTasks, err := r.sdk.Tasks().Create(ctx, []models.Task{task})
	if err != nil {
		return nil, err
	}

	if len(createdTasks) == 0 {
		return nil, fmt.Errorf("failed to create task: empty response")
	}

	return &createdTasks[0], nil
}
