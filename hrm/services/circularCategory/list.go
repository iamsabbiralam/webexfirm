package circularCategory

import (
	"context"

	cc "practice/webex/gunk/v1/circularCategory"
	"practice/webex/hrm/storage"

	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) ListCircularCategory(ctx context.Context, req *cc.ListCircularCategoryRequest) (*cc.ListCircularCategoryResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.CircularCategory.ListCircularCategory")
	dbPrm := storage.CircularCategory{
		SearchTerm: req.SearchTerm,
		Offset:     req.Offset,
		Limit:      req.Limit,
		Status:     int32(req.Status),
	}

	ccList, err := h.ccst.ListCircularCategory(ctx, dbPrm)
	if err != nil {
		errMsg := "failed to list circular category"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	list := make([]*cc.CircularCategory, len(ccList))
	for i, c := range ccList {
		list[i] = &cc.CircularCategory{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			Status:      cc.Status(c.Status),
			Position:    int64(c.Position),
			CreatedAt:   tspb.New(c.CreatedAt),
			CreatedBy:   c.CreatedBy,
			UpdatedAt:   tspb.New(c.UpdatedAt),
			UpdatedBy:   c.UpdatedBy,
		}
	}

	var total int32
	if len(ccList) > 0 {
		total = int32(ccList[0].Count)
	}

	return &cc.ListCircularCategoryResponse{
		CircularCategory: list,
		Total:            total,
	}, nil
}
