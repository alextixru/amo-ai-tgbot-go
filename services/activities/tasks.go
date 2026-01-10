package activities

import (
	"context"
	"time"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
)

func (s *service) ListTasks(ctx context.Context, parent *gkitmodels.ParentEntity, filter *gkitmodels.ActivityFilter) ([]*models.Task, error) {
	f := filters.NewTasksFilter()
	if filter != nil && filter.Limit > 0 {
		f.SetLimit(filter.Limit)
	} else {
		f.SetLimit(50)
	}

	if parent != nil {
		f.SetEntityType(parent.Type)
		f.SetEntityIDs([]int{parent.ID})
	}

	if filter != nil {
		if len(filter.ResponsibleUserID) > 0 {
			f.SetResponsibleUserIDs(filter.ResponsibleUserID)
		}
		if filter.IsCompleted != nil {
			f.SetIsCompleted(*filter.IsCompleted)
		}
		if len(filter.TaskTypeID) > 0 {
			// SDK expects int for SetTaskTypeID, but multiple IDs might not be supported by SDK filter directly
			// Based on filter/tasks.go: f.TaskTypeID = &id
			// We take the first one for now or we could use multi-task-type if API supports it (it usually doesn't in v4 filter[task_type_id])
			if len(filter.TaskTypeID) > 0 {
				f.SetTaskTypeID(filter.TaskTypeID[0])
			}
		}
	}

	// Always sort by deadline to allow effective client-side filtering
	f.SetOrder("complete_till", "asc")

	tasks, _, err := s.sdk.Tasks().Get(ctx, f)
	if err != nil {
		return nil, err
	}

	if filter == nil || filter.DateRange == "" {
		return tasks, nil
	}

	// Client-side filtering for DateRange
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
		switch filter.DateRange {
		case "today":
			if taskTime.After(todayStart) && taskTime.Before(todayEnd) {
				include = true
			}
		case "tomorrow":
			tomorrowStart := todayStart.Add(24 * time.Hour)
			tomorrowEnd := tomorrowStart.Add(24 * time.Hour)
			if taskTime.After(tomorrowStart) && taskTime.Before(tomorrowEnd) {
				include = true
			}
		case "overdue":
			if taskTime.Before(todayStart) && (task.IsCompleted == false) {
				include = true
			}
		case "future":
			if taskTime.After(todayEnd) {
				include = true
			}
		default:
			include = true
		}

		if include {
			filtered = append(filtered, task)
		}
	}

	return filtered, nil
}

func (s *service) GetTask(ctx context.Context, id int) (*models.Task, error) {
	return s.sdk.Tasks().GetOne(ctx, id)
}

func (s *service) CreateTask(ctx context.Context, parent gkitmodels.ParentEntity, data *gkitmodels.ActivityData) (*models.Task, error) {
	task := &models.Task{
		Text:       data.Text,
		EntityID:   parent.ID,
		EntityType: parent.Type,
	}
	if data.CompleteTillAt > 0 {
		task.CompleteTill = &data.CompleteTillAt
	} else if data.CompleteTill > 0 {
		task.CompleteTill = &data.CompleteTill
	}
	if data.TaskTypeID > 0 {
		task.TaskTypeID = data.TaskTypeID
	} else if data.TaskType > 0 {
		task.TaskTypeID = data.TaskType
	}
	if data.ResponsibleUserID > 0 {
		task.ResponsibleUserID = data.ResponsibleUserID
	}
	tasks, _, err := s.sdk.Tasks().Create(ctx, []*models.Task{task})
	if err != nil {
		return nil, err
	}
	if len(tasks) > 0 {
		return tasks[0], nil
	}
	return nil, nil
}

func (s *service) UpdateTask(ctx context.Context, id int, data *gkitmodels.ActivityData) (*models.Task, error) {
	task := &models.Task{
		BaseModel: models.BaseModel{ID: id},
	}
	if data.Text != "" {
		task.Text = data.Text
	}
	if data.CompleteTillAt > 0 {
		task.CompleteTill = &data.CompleteTillAt
	} else if data.CompleteTill > 0 {
		task.CompleteTill = &data.CompleteTill
	}
	tasks, _, err := s.sdk.Tasks().Update(ctx, []*models.Task{task})
	if err != nil {
		return nil, err
	}
	if len(tasks) > 0 {
		return tasks[0], nil
	}
	return nil, nil
}

func (s *service) CompleteTask(ctx context.Context, id int, resultText string) (*models.Task, error) {
	return s.sdk.Tasks().Complete(ctx, id, resultText)
}
