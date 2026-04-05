package customers

import (
	"context"

	"github.com/alextixru/amocrm-sdk-go/core/filters"
	"github.com/alextixru/amocrm-sdk-go/core/models"
)

func (s *service) ListSegments(ctx context.Context, page, limit int) (*SegmentsListOutput, error) {
	f := filters.NewSegmentsFilter()
	if page > 0 {
		f.Page = page
	}
	if limit > 0 {
		f.Limit = limit
	}
	res, meta, err := s.sdk.Segments().Get(ctx, f)
	if err != nil {
		return nil, err
	}
	out := &SegmentsListOutput{
		Segments: make([]*SegmentOutput, 0, len(res)),
	}
	for _, seg := range res {
		out.Segments = append(out.Segments, segmentToOutput(seg))
	}
	if meta != nil {
		out.HasMore = meta.HasMore
	}
	return out, nil
}

func (s *service) GetSegment(ctx context.Context, id int) (*SegmentOutput, error) {
	seg, err := s.sdk.Segments().GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return segmentToOutput(seg), nil
}

func (s *service) CreateSegments(ctx context.Context, names []string) (*SegmentsListOutput, error) {
	segments := make([]*models.CustomerSegment, 0, len(names))
	for _, name := range names {
		segments = append(segments, &models.CustomerSegment{Name: name})
	}
	res, meta, err := s.sdk.Segments().Create(ctx, segments)
	if err != nil {
		return nil, err
	}
	out := &SegmentsListOutput{
		Segments: make([]*SegmentOutput, 0, len(res)),
	}
	for _, seg := range res {
		out.Segments = append(out.Segments, segmentToOutput(seg))
	}
	if meta != nil {
		out.HasMore = meta.HasMore
	}
	return out, nil
}

func (s *service) DeleteSegment(ctx context.Context, id int) error {
	return s.sdk.Segments().Delete(ctx, id)
}

// segmentToOutput конвертирует *models.CustomerSegment в *SegmentOutput.
func segmentToOutput(seg *models.CustomerSegment) *SegmentOutput {
	if seg == nil {
		return nil
	}
	out := &SegmentOutput{
		ID:    seg.ID,
		Name:  seg.Name,
		Color: seg.Color,
	}
	if seg.CustomersCount != nil {
		out.CustomersCount = *seg.CustomersCount
	}
	return out
}
