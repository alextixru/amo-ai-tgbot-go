package customers

import (
	"context"
	"fmt"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
	gkitmodels "github.com/tihn/amo-ai-tgbot-go/internal/models/tools"
)

// defaultWith набор связанных сущностей, запрашиваемых по умолчанию при get.
var defaultWith = []string{"contacts", "companies", "segments"}

func (s *service) ListCustomers(ctx context.Context, filter *gkitmodels.CustomerFilter, with []string) (*CustomersListOutput, error) {
	sdkFilter := filters.NewCustomersFilter()

	if filter != nil {
		if filter.NextDateFrom != "" || filter.NextDateTo != "" {
			return nil, fmt.Errorf("фильтр next_date_from/next_date_to не поддерживается текущей версией API amoCRM")
		}
		if filter.Page > 0 {
			sdkFilter.SetPage(filter.Page)
		}
		if filter.Limit > 0 {
			sdkFilter.SetLimit(filter.Limit)
		}
		if filter.Query != "" {
			sdkFilter.SetQuery(filter.Query)
		}
		if len(filter.IDs) > 0 {
			sdkFilter.SetIDs(filter.IDs)
		}
		if len(filter.Names) > 0 {
			sdkFilter.SetNames(filter.Names)
		}
		if len(filter.ResponsibleUserNames) > 0 {
			ids, err := s.resolveUserNames(filter.ResponsibleUserNames)
			if err != nil {
				return nil, err
			}
			sdkFilter.SetResponsibleUserIDs(ids)
		}
		if len(filter.StatusNames) > 0 {
			ids, err := s.resolveStatusNames(filter.StatusNames)
			if err != nil {
				return nil, err
			}
			sdkFilter.SetStatusIDs(ids)
		}
	}

	if len(with) > 0 {
		sdkFilter.With = with
	}

	customers, meta, err := s.sdk.Customers().Get(ctx, sdkFilter)
	if err != nil {
		return nil, err
	}

	out := &CustomersListOutput{
		Customers: make([]*CustomerOutput, 0, len(customers)),
	}
	for i := range customers {
		out.Customers = append(out.Customers, s.enrichCustomer(&customers[i]))
	}
	if meta != nil {
		out.HasMore = meta.HasMore
	}
	return out, nil
}

func (s *service) GetCustomer(ctx context.Context, id int, with []string) (*CustomerOutput, error) {
	sdkFilter := filters.NewCustomersFilter()
	sdkFilter.SetIDs([]int{id})

	// По умолчанию запрашиваем связанные данные
	if len(with) > 0 {
		sdkFilter.With = with
	} else {
		sdkFilter.With = defaultWith
	}

	customers, _, err := s.sdk.Customers().Get(ctx, sdkFilter)
	if err != nil {
		return nil, err
	}
	if len(customers) == 0 {
		return nil, fmt.Errorf("покупатель с ID %d не найден", id)
	}
	return s.enrichCustomer(&customers[0]), nil
}

func (s *service) CreateCustomers(ctx context.Context, data []*gkitmodels.CustomerData) ([]*CustomerOutput, error) {
	customers := make([]models.Customer, 0, len(data))
	for _, d := range data {
		c, err := s.mapCustomerData(d, 0)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}

	res, _, err := s.sdk.Customers().Create(ctx, customers)
	if err != nil {
		return nil, err
	}

	out := make([]*CustomerOutput, 0, len(res))
	for i := range res {
		out = append(out, s.enrichCustomer(&res[i]))
	}
	return out, nil
}

func (s *service) UpdateCustomers(ctx context.Context, id int, data []*gkitmodels.CustomerData) ([]*CustomerOutput, error) {
	customers := make([]models.Customer, 0, len(data))
	for _, d := range data {
		c, err := s.mapCustomerData(d, id)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}

	res, _, err := s.sdk.Customers().Update(ctx, customers)
	if err != nil {
		return nil, err
	}

	out := make([]*CustomerOutput, 0, len(res))
	for i := range res {
		out = append(out, s.enrichCustomer(&res[i]))
	}
	return out, nil
}

func (s *service) DeleteCustomer(ctx context.Context, id int) error {
	return s.sdk.Customers().Delete(ctx, id)
}

func (s *service) LinkCustomer(ctx context.Context, customerID int, entityType string, entityID int) error {
	link := models.EntityLink{
		ToEntityType: entityType,
		ToEntityID:   entityID,
	}
	_, err := s.sdk.Customers().Link(ctx, customerID, []models.EntityLink{link})
	return err
}

// mapCustomerData конвертирует CustomerData в SDK Customer, резолвя имена в ID.
func (s *service) mapCustomerData(d *gkitmodels.CustomerData, id int) (models.Customer, error) {
	c := models.Customer{
		Name:        d.Name,
		NextPrice:   d.NextPrice,
		Periodicity: d.Periodicity,
	}
	c.ID = id

	if d.NextDate != "" {
		ts, err := parseISO(d.NextDate)
		if err != nil {
			return models.Customer{}, err
		}
		c.NextDate = ts
	}

	if d.ResponsibleUserName != "" {
		uid, err := s.resolveUserName(d.ResponsibleUserName)
		if err != nil {
			return models.Customer{}, err
		}
		c.ResponsibleUserID = uid
	}

	if d.StatusName != "" {
		sid, err := s.resolveStatusName(d.StatusName)
		if err != nil {
			return models.Customer{}, err
		}
		c.StatusID = sid
	}

	if len(d.CustomFieldsValues) > 0 {
		c.CustomFieldsValues = mapCustomerFieldValues(d.CustomFieldsValues)
	}

	for _, tag := range d.TagsToAdd {
		c.TagsToAdd = append(c.TagsToAdd, models.Tag{Name: tag})
	}
	for _, tag := range d.TagsToDelete {
		c.TagsToDelete = append(c.TagsToDelete, models.Tag{Name: tag})
	}

	return c, nil
}

// mapCustomerFieldValues конвертирует []CustomerFieldValue в []models.CustomFieldValue.
func mapCustomerFieldValues(cfv []gkitmodels.CustomerFieldValue) []models.CustomFieldValue {
	result := make([]models.CustomFieldValue, 0, len(cfv))
	for _, f := range cfv {
		elem := models.FieldValueElement{Value: f.Value}
		if f.EnumCode != "" {
			elem.EnumCode = f.EnumCode
		}
		result = append(result, models.CustomFieldValue{
			FieldCode: f.FieldCode,
			Values:    []models.FieldValueElement{elem},
		})
	}
	return result
}

// enrichCustomer конвертирует SDK Customer в CustomerOutput с именами вместо ID.
func (s *service) enrichCustomer(c *models.Customer) *CustomerOutput {
	if c == nil {
		return nil
	}
	out := &CustomerOutput{
		ID:                  c.ID,
		Name:                c.Name,
		NextPrice:           c.NextPrice,
		NextDate:            toISO(c.NextDate),
		StatusName:          s.resolveStatusID(c.StatusID),
		Periodicity:         c.Periodicity,
		ResponsibleUserName: s.resolveUserID(c.ResponsibleUserID),
		CreatedByName:       s.resolveUserID(c.CreatedBy),
		UpdatedByName:       s.resolveUserID(c.UpdatedBy),
		CreatedAt:           toISO(c.CreatedAt),
		UpdatedAt:           toISO(c.UpdatedAt),
		IsDeleted:           c.IsDeleted,
		Ltv:                 c.Ltv,
		PurchasesCount:      c.PurchasesCount,
		AverageCheck:        c.AverageCheck,
	}

	if c.Embedded != nil {
		for _, t := range c.Embedded.Tags {
			out.Tags = append(out.Tags, t.Name)
		}
		for _, seg := range c.Embedded.Segments {
			out.Segments = append(out.Segments, CustomerSegmentBrief{ID: seg.ID, Name: seg.Name})
		}
		for _, ct := range c.Embedded.Contacts {
			if ct != nil {
				out.Contacts = append(out.Contacts, CustomerEntityBrief{ID: ct.ID, Name: ct.Name})
			}
		}
		for _, co := range c.Embedded.Companies {
			if co != nil {
				out.Companies = append(out.Companies, CustomerEntityBrief{ID: co.ID, Name: co.Name})
			}
		}
	}

	return out
}
