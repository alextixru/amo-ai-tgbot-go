package activities

import (
	"context"
	"fmt"
	"time"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
	"github.com/tihn/amo-ai-tgbot-go/internal/utils"
)

// taskTypeNameToID конвертирует строковое имя типа задачи в SDK-константу.
func taskTypeNameToID(name string) int {
	switch name {
	case "follow_up":
		return int(models.TaskTypeFollowUp) // 1
	case "meeting":
		return int(models.TaskTypeMeeting) // 2
	default:
		return 0
	}
}

// taskTypeIDToName конвертирует SDK-константу в читаемое имя.
func taskTypeIDToName(id int) string {
	switch models.TaskTypeID(id) {
	case models.TaskTypeFollowUp:
		return "follow_up"
	case models.TaskTypeMeeting:
		return "meeting"
	default:
		if id > 0 {
			return fmt.Sprintf("type_%d", id)
		}
		return ""
	}
}

// toISO конвертирует Unix timestamp в ISO 8601. Возвращает "" если ts==0.
func toISO(ts int64) string {
	if ts == 0 {
		return ""
	}
	return time.Unix(ts, 0).UTC().Format(time.RFC3339)
}

// convertTask конвертирует SDK-задачу в TaskOutput.
func (s *service) convertTask(t *models.Task) *TaskOutput {
	if t == nil {
		return nil
	}
	out := &TaskOutput{
		ID:                  t.ID,
		Text:                t.Text,
		EntityID:            t.EntityID,
		EntityType:          t.EntityType,
		TaskType:            taskTypeIDToName(t.TaskTypeID),
		IsCompleted:         t.IsCompleted,
		ResponsibleUserName: s.resolveUserID(t.ResponsibleUserID),
		CreatedByName:       s.resolveUserID(t.CreatedBy),
		UpdatedByName:       s.resolveUserID(t.UpdatedBy),
		CreatedAt:           toISO(t.CreatedAt),
		UpdatedAt:           toISO(t.UpdatedAt),
	}
	if t.CompleteTill != nil {
		out.Deadline = toISO(*t.CompleteTill)
	}
	if t.Result != nil && t.Result.Text != "" {
		out.Result = &TaskResult{Text: t.Result.Text}
	}
	return out
}

func (s *service) ListTasks(ctx context.Context, parent *gkitmodels.ParentEntity, filter *gkitmodels.TasksFilter, with []string) (*TasksListOutput, error) {
	f := filters.NewTasksFilter()
	if filter != nil && filter.Limit > 0 {
		f.SetLimit(filter.Limit)
	} else {
		f.SetLimit(50)
	}

	if filter != nil {
		if filter.Page > 0 {
			f.SetPage(filter.Page)
		}
		if filter.Order != "" {
			dir := "asc"
			if filter.OrderDir != "" {
				dir = filter.OrderDir
			}
			f.SetOrder(filter.Order, dir)
		} else {
			f.SetOrder("complete_till", "asc")
		}
		if len(filter.IDs) > 0 {
			f.SetIDs(filter.IDs)
		}
		if filter.UpdatedAt != nil {
			from := int(*filter.UpdatedAt)
			var toPtr *int
			if filter.UpdatedAtTo != nil {
				to := int(*filter.UpdatedAtTo)
				toPtr = &to
			}
			f.SetUpdatedAt(&from, toPtr)
		}
		// Resolve responsible user names → IDs
		if len(filter.ResponsibleUserNames) > 0 {
			ids, err := s.resolveUserNames(filter.ResponsibleUserNames)
			if err != nil {
				return nil, err
			}
			f.SetResponsibleUserIDs(ids)
		}
		// Resolve created_by names → IDs
		if len(filter.CreatedByNames) > 0 {
			ids, err := s.resolveUserNames(filter.CreatedByNames)
			if err != nil {
				return nil, err
			}
			f.SetCreatedBy(ids)
		}
		if filter.IsCompleted != nil {
			f.SetIsCompleted(*filter.IsCompleted)
		}
		if filter.TaskType != "" {
			typeID := taskTypeNameToID(filter.TaskType)
			if typeID > 0 {
				f.SetTaskTypeID(typeID)
			}
		}
	} else {
		f.SetOrder("complete_till", "asc")
	}

	if parent != nil {
		f.SetEntityType(parent.Type)
		f.SetEntityIDs([]int{parent.ID})
	}

	if len(with) > 0 {
		f.With = with
	}

	tasks, meta, err := s.sdk.Tasks().Get(ctx, f)
	if err != nil {
		return nil, err
	}

	// Client-side filtering for DateRange
	if filter != nil && filter.DateRange != "" {
		tasks = filterByDateRange(tasks, filter.DateRange)
	}

	out := &TasksListOutput{
		Tasks: make([]*TaskOutput, 0, len(tasks)),
	}
	for _, t := range tasks {
		out.Tasks = append(out.Tasks, s.convertTask(t))
	}
	if meta != nil {
		out.PageMeta = PageMeta{HasMore: meta.HasMore, Total: meta.TotalItems}
	}
	return out, nil
}

func (s *service) GetTask(ctx context.Context, id int, with []string) (*TaskOutput, error) {
	f := filters.NewTasksFilter()
	f.SetIDs([]int{id})
	if len(with) > 0 {
		f.With = with
	}
	tasks, _, err := s.sdk.Tasks().Get(ctx, f)
	if err != nil {
		return nil, err
	}
	if len(tasks) > 0 {
		return s.convertTask(tasks[0]), nil
	}
	return nil, nil
}

func (s *service) CreateTask(ctx context.Context, parent gkitmodels.ParentEntity, data *gkitmodels.TaskData) (*TaskOutput, error) {
	result, err := s.CreateTasks(ctx, parent, []gkitmodels.TaskData{*data})
	if err != nil {
		return nil, err
	}
	if result != nil && len(result.Tasks) > 0 {
		return result.Tasks[0], nil
	}
	return nil, nil
}

func (s *service) CreateTasks(ctx context.Context, parent gkitmodels.ParentEntity, data []gkitmodels.TaskData) (*TasksListOutput, error) {
	tasks := make([]*models.Task, len(data))
	for i, d := range data {
		task := &models.Task{
			Text:       d.Text,
			EntityID:   parent.ID,
			EntityType: parent.Type,
		}
		if d.Deadline != "" {
			if ts, err := utils.ParseHumanDeadline(d.Deadline); err == nil && ts > 0 {
				task.CompleteTill = &ts
			}
		}
		if d.TaskType != "" {
			typeID := taskTypeNameToID(d.TaskType)
			if typeID > 0 {
				task.TaskTypeID = typeID
			}
		}
		if d.ResponsibleUserName != "" {
			uid, err := s.resolveUserName(d.ResponsibleUserName)
			if err != nil {
				return nil, err
			}
			task.ResponsibleUserID = uid
		}
		tasks[i] = task
	}

	result, _, err := s.sdk.Tasks().Create(ctx, tasks)
	if err != nil {
		return nil, err
	}
	out := &TasksListOutput{Tasks: make([]*TaskOutput, 0, len(result))}
	for _, t := range result {
		out.Tasks = append(out.Tasks, s.convertTask(t))
	}
	return out, nil
}

func (s *service) UpdateTask(ctx context.Context, id int, data *gkitmodels.TaskData) (*TaskOutput, error) {
	task := &models.Task{
		BaseModel: models.BaseModel{ID: id},
	}
	if data.Text != "" {
		task.Text = data.Text
	}
	if data.Deadline != "" {
		if ts, err := utils.ParseHumanDeadline(data.Deadline); err == nil && ts > 0 {
			task.CompleteTill = &ts
		}
	}
	if data.TaskType != "" {
		typeID := taskTypeNameToID(data.TaskType)
		if typeID > 0 {
			task.TaskTypeID = typeID
		}
	}
	if data.ResponsibleUserName != "" {
		uid, err := s.resolveUserName(data.ResponsibleUserName)
		if err != nil {
			return nil, err
		}
		task.ResponsibleUserID = uid
	}
	tasks, _, err := s.sdk.Tasks().Update(ctx, []*models.Task{task})
	if err != nil {
		return nil, err
	}
	if len(tasks) > 0 {
		return s.convertTask(tasks[0]), nil
	}
	return nil, nil
}

func (s *service) CompleteTask(ctx context.Context, id int, resultText string) (*TaskOutput, error) {
	t, err := s.sdk.Tasks().Complete(ctx, id, resultText)
	if err != nil {
		return nil, err
	}
	return s.convertTask(t), nil
}

// filterByDateRange фильтрует задачи по временному диапазону на стороне клиента.
func filterByDateRange(tasks []*models.Task, dateRange string) []*models.Task {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	var filtered []*models.Task
	for _, task := range tasks {
		if task.CompleteTill == nil {
			continue
		}
		taskTime := time.Unix(*task.CompleteTill, 0)

		include := false
		switch dateRange {
		case "today":
			include = taskTime.After(todayStart) && taskTime.Before(todayEnd)
		case "tomorrow":
			tomorrowStart := todayStart.Add(24 * time.Hour)
			tomorrowEnd := tomorrowStart.Add(24 * time.Hour)
			include = taskTime.After(tomorrowStart) && taskTime.Before(tomorrowEnd)
		case "overdue":
			include = taskTime.Before(todayStart) && !task.IsCompleted
		case "this_week":
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			mondayStart := todayStart.AddDate(0, 0, -(weekday - 1))
			sundayEnd := mondayStart.AddDate(0, 0, 7)
			include = taskTime.After(mondayStart) && taskTime.Before(sundayEnd)
		case "next_week":
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			nextMondayStart := todayStart.AddDate(0, 0, -(weekday-1)+7)
			nextSundayEnd := nextMondayStart.AddDate(0, 0, 7)
			include = taskTime.After(nextMondayStart) && taskTime.Before(nextSundayEnd)
		default:
			include = true
		}

		if include {
			filtered = append(filtered, task)
		}
	}
	return filtered
}
