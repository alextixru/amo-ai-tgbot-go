package entities

import (
	"context"
	"fmt"

	sdkfilters "github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

func (s *service) SearchLeads(ctx context.Context, filter *gkitmodels.EntitiesFilter, with []string) (*SearchResult, error) {
	f := sdkfilters.NewLeadsFilter()
	f.SetLimit(50)
	f.SetPage(1)

	if filter != nil {
		if filter.Query != "" {
			f.SetQuery(filter.Query)
		}
		if filter.Limit > 0 {
			f.SetLimit(filter.Limit)
		}
		if filter.Page > 0 {
			f.SetPage(filter.Page)
		}
		if len(filter.IDs) > 0 {
			f.SetIDs(filter.IDs)
		}

		// Ответственные: имена → ID
		if len(filter.ResponsibleUserNames) > 0 {
			ids, err := s.resolveUserIDs(filter.ResponsibleUserNames)
			if err != nil {
				return nil, err
			}
			f.SetResponsibleUserIDs(ids)
		}

		// created_by: имена → ID
		if len(filter.CreatedByNames) > 0 {
			ids, err := s.resolveUserIDs(filter.CreatedByNames)
			if err != nil {
				return nil, err
			}
			f.SetCreatedBy(ids)
		}

		// updated_by: имена → ID
		if len(filter.UpdatedByNames) > 0 {
			ids, err := s.resolveUserIDs(filter.UpdatedByNames)
			if err != nil {
				return nil, err
			}
			f.SetUpdatedBy(ids)
		}

		// Воронки: имена → ID
		if len(filter.PipelineNames) > 0 {
			ids := make([]int, 0, len(filter.PipelineNames))
			for _, name := range filter.PipelineNames {
				id, err := s.resolvePipelineID(name)
				if err != nil {
					return nil, err
				}
				ids = append(ids, id)
			}
			f.SetPipelineIDs(ids)
		}

		// Статусы: пары {pipeline_name, status_name} → []LeadStatusFilter
		if len(filter.Statuses) > 0 {
			statusFilters := make([]sdkfilters.LeadStatusFilter, 0, len(filter.Statuses))
			for _, sp := range filter.Statuses {
				pipelineID, statusID, err := s.resolveStatusID(sp.PipelineName, sp.StatusName)
				if err != nil {
					return nil, err
				}
				statusFilters = append(statusFilters, sdkfilters.LeadStatusFilter{
					PipelineID: pipelineID,
					StatusID:   statusID,
				})
			}
			f.SetStatuses(statusFilters)
		}

		// Цена
		if filter.PriceFrom > 0 || filter.PriceTo > 0 {
			var from, to *int
			if filter.PriceFrom > 0 {
				v := filter.PriceFrom
				from = &v
			}
			if filter.PriceTo > 0 {
				v := filter.PriceTo
				to = &v
			}
			f.SetPrice(from, to)
		}

		// Даты (ISO-8601 → Unix)
		if filter.CreatedAtFrom != "" || filter.CreatedAtTo != "" {
			from := parseISO(filter.CreatedAtFrom)
			to := parseISO(filter.CreatedAtTo)
			f.SetCreatedAt(intPtrOrNil(from), intPtrOrNil(to))
		}
		if filter.UpdatedAtFrom != "" || filter.UpdatedAtTo != "" {
			from := parseISO(filter.UpdatedAtFrom)
			to := parseISO(filter.UpdatedAtTo)
			f.SetUpdatedAt(intPtrOrNil(from), intPtrOrNil(to))
		}
		if filter.ClosedAtFrom != "" || filter.ClosedAtTo != "" {
			from := parseISO(filter.ClosedAtFrom)
			to := parseISO(filter.ClosedAtTo)
			f.SetClosedAt(intPtrOrNil(from), intPtrOrNil(to))
		}

		// Кастомные поля: field_code → field_id
		if len(filter.CustomFieldsValues) > 0 {
			cfMap := buildCustomFieldsFilter(filter.CustomFieldsValues, s.customFieldsLeads)
			if len(cfMap) > 0 {
				f.SetCustomFieldsValues(cfMap)
			}
		}

		if len(filter.With) > 0 {
			with = append(with, filter.With...)
		}
	}

	if len(with) > 0 {
		f.With = with
	}

	leads, meta, err := s.sdk.Leads().Get(ctx, f)
	if err != nil {
		return nil, err
	}

	items := make([]*EntityResult, 0, len(leads))
	for _, lead := range leads {
		items = append(items, s.leadToResult(lead))
	}

	result := &SearchResult{
		Items: items,
		Page:  1,
	}
	if meta != nil {
		result.HasMore = meta.HasMore
		result.Page = meta.Page
	}
	return result, nil
}

func (s *service) GetLead(ctx context.Context, id int, with []string) (*EntityResult, error) {
	// Дефолтные with для обогащения ответа
	defaults := []string{"loss_reason", "source"}
	for _, d := range defaults {
		if !containsStr(with, d) {
			with = append(with, d)
		}
	}

	f := sdkfilters.NewLeadsFilter()
	f.SetIDs([]int{id})
	if len(with) > 0 {
		f.With = with
	}

	leads, _, err := s.sdk.Leads().Get(ctx, f)
	if err != nil {
		return nil, err
	}
	if len(leads) == 0 {
		return nil, nil
	}
	return s.leadToResult(leads[0]), nil
}

func (s *service) CreateLead(ctx context.Context, data *gkitmodels.EntityData) (*EntityResult, error) {
	lead, err := s.mapToLead(data)
	if err != nil {
		return nil, err
	}
	created, err := s.sdk.Leads().CreateOne(ctx, lead)
	if err != nil {
		return nil, err
	}
	return s.leadToResult(created), nil
}

func (s *service) CreateLeads(ctx context.Context, dataList []gkitmodels.EntityData) ([]*EntityResult, error) {
	leads := make([]*models.Lead, 0, len(dataList))
	for i := range dataList {
		lead, err := s.mapToLead(&dataList[i])
		if err != nil {
			return nil, fmt.Errorf("item %d: %w", i, err)
		}
		leads = append(leads, lead)
	}
	created, _, err := s.sdk.Leads().Create(ctx, leads)
	if err != nil {
		return nil, err
	}
	results := make([]*EntityResult, 0, len(created))
	for _, l := range created {
		results = append(results, s.leadToResult(l))
	}
	return results, nil
}

func (s *service) UpdateLead(ctx context.Context, id int, data *gkitmodels.EntityData) (*EntityResult, error) {
	lead, err := s.mapToLead(data)
	if err != nil {
		return nil, err
	}
	lead.ID = id
	updated, err := s.sdk.Leads().UpdateOne(ctx, lead)
	if err != nil {
		return nil, err
	}
	return s.leadToResult(updated), nil
}

func (s *service) UpdateLeads(ctx context.Context, dataList []gkitmodels.EntityData) ([]*EntityResult, error) {
	leads := make([]*models.Lead, 0, len(dataList))
	for i := range dataList {
		lead, err := s.mapToLead(&dataList[i])
		if err != nil {
			return nil, fmt.Errorf("item %d: %w", i, err)
		}
		if lead.ID == 0 && dataList[i].ID > 0 {
			lead.ID = dataList[i].ID
		}
		leads = append(leads, lead)
	}
	updated, _, err := s.sdk.Leads().Update(ctx, leads)
	if err != nil {
		return nil, err
	}
	results := make([]*EntityResult, 0, len(updated))
	for _, l := range updated {
		results = append(results, s.leadToResult(l))
	}
	return results, nil
}

func (s *service) SyncLead(ctx context.Context, id int, data *gkitmodels.EntityData) (*EntityResult, error) {
	lead, err := s.mapToLead(data)
	if err != nil {
		return nil, err
	}
	if id > 0 {
		lead.ID = id
	}
	synced, err := s.sdk.Leads().SyncOne(ctx, lead, []string{"contacts", "companies"})
	if err != nil {
		return nil, err
	}
	return s.leadToResult(synced), nil
}

func (s *service) LinkLead(ctx context.Context, leadID int, target *gkitmodels.LinkTarget) (*LinkResult, error) {
	link := models.EntityLink{
		ToEntityID:   target.ID,
		ToEntityType: target.Type,
	}
	_, err := s.sdk.Leads().Link(ctx, leadID, []models.EntityLink{link})
	if err != nil {
		return nil, err
	}
	return &LinkResult{
		Success: true,
		Message: fmt.Sprintf("Сделка %d связана с %s %d", leadID, target.Type, target.ID),
	}, nil
}

func (s *service) UnlinkLead(ctx context.Context, leadID int, target *gkitmodels.LinkTarget) (*LinkResult, error) {
	link := models.EntityLink{
		ToEntityID:   target.ID,
		ToEntityType: target.Type,
	}
	if err := s.sdk.Leads().Unlink(ctx, leadID, []models.EntityLink{link}); err != nil {
		return nil, err
	}
	return &LinkResult{
		Success: true,
		Message: fmt.Sprintf("Сделка %d отвязана от %s %d", leadID, target.Type, target.ID),
	}, nil
}

// mapToLead конвертирует EntityData в SDK Lead, резолвя имена в ID.
func (s *service) mapToLead(data *gkitmodels.EntityData) (*models.Lead, error) {
	lead := &models.Lead{
		Name:  data.Name,
		Price: data.Price,
	}
	if data.ID > 0 {
		lead.ID = data.ID
	}

	// Воронка и статус — имена → ID
	if data.PipelineName != "" {
		id, err := s.resolvePipelineID(data.PipelineName)
		if err != nil {
			return nil, err
		}
		lead.PipelineID = id
	}
	if data.StatusName != "" {
		_, statusID, err := s.resolveStatusID(data.PipelineName, data.StatusName)
		if err != nil {
			return nil, err
		}
		lead.StatusID = statusID
	}

	// Ответственный — имя → ID
	if data.ResponsibleUserName != "" {
		id, err := s.resolveUserID(data.ResponsibleUserName)
		if err != nil {
			return nil, err
		}
		lead.ResponsibleUserID = id
	}

	// Причина отказа — имя → ID
	if data.LossReasonName != "" {
		id, err := s.resolveLossReasonID(data.LossReasonName)
		if err != nil {
			return nil, err
		}
		lead.LossReasonID = &id
	}

	lead.CustomFieldsValues = mapCustomFieldsValues(data.CustomFieldsValues)

	if len(data.Tags) > 0 || len(data.EmbeddedContacts) > 0 || len(data.EmbeddedCompanies) > 0 {
		lead.Embedded = &models.LeadEmbedded{}
		if len(data.Tags) > 0 {
			lead.Embedded.Tags = mapTags(data.Tags)
		}
		if len(data.EmbeddedContacts) > 0 {
			lead.Embedded.Contacts = make([]*models.Contact, len(data.EmbeddedContacts))
			for i, cid := range data.EmbeddedContacts {
				lead.Embedded.Contacts[i] = &models.Contact{}
				lead.Embedded.Contacts[i].ID = cid
			}
		}
		if len(data.EmbeddedCompanies) > 0 {
			lead.Embedded.Companies = make([]*models.Company, len(data.EmbeddedCompanies))
			for i, cid := range data.EmbeddedCompanies {
				lead.Embedded.Companies[i] = &models.Company{}
				lead.Embedded.Companies[i].ID = cid
			}
		}
	}

	return lead, nil
}

// intPtrOrNil возвращает указатель на int, или nil если значение 0.
func intPtrOrNil(v int) *int {
	if v == 0 {
		return nil
	}
	return &v
}

// containsStr проверяет наличие строки в срезе.
func containsStr(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
