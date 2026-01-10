package activities

import (
	"context"
	"time"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/models"
	"github.com/tihn/amo-ai-tgbot-go/utils"
)

func (s *service) ListTasks(ctx context.Context, parent *gkitmodels.ParentEntity, filter *gkitmodels.TasksFilter) ([]*models.Task, error) {
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
		if filter.TaskTypeID > 0 {
			f.SetTaskTypeID(filter.TaskTypeID)
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
		case "this_week":
			// Monday to Sunday
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			mondayStart := todayStart.AddDate(0, 0, -(weekday - 1))
			sundayEnd := mondayStart.AddDate(0, 0, 7)
			if taskTime.After(mondayStart) && taskTime.Before(sundayEnd) {
				include = true
			}
		case "next_week":
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			nextMondayStart := todayStart.AddDate(0, 0, -(weekday-1)+7)
			nextSundayEnd := nextMondayStart.AddDate(0, 0, 7)
			if taskTime.After(nextMondayStart) && taskTime.Before(nextSundayEnd) {
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

func (s *service) CreateTask(ctx context.Context, parent gkitmodels.ParentEntity, data *gkitmodels.TaskData) (*models.Task, error) {
	task := &models.Task{
		Text:       data.Text,
		EntityID:   parent.ID,
		EntityType: parent.Type,
	}

	if data.Deadline != "" {
		if ts, err := utils.ParseHumanDeadline(data.Deadline); err == nil && ts > 0 {
			task.CompleteTill = &ts
		}
	}

	if data.TaskTypeID > 0 {
		task.TaskTypeID = data.TaskTypeID
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

func (s *service) UpdateTask(ctx context.Context, id int, data *gkitmodels.TaskData) (*models.Task, error) {
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
