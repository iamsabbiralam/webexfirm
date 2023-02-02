package circularCategory

import (
	"context"

	cc "practice/webex/gunk/v1/circularCategory"
	"practice/webex/hrm/storage"

	"google.golang.org/grpc/status"
	"practice/webex/serviceutil/logging"
)

func (h *Handler) CreateCircularCategory(ctx context.Context, req *cc.CreateCircularCategoryRequest) (*cc.CreateCircularCategoryResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.circular-category.CreateCircularCategory")
	dbPrm := formatCircularCategory(req)
	id, err := h.ccst.CreateCircularCategory(ctx, dbPrm)
	if err != nil {
		errMsg := "failed to create circular category"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	return &cc.CreateCircularCategoryResponse{
		ID: id,
	}, nil
}

func formatCircularCategory(req *cc.CreateCircularCategoryRequest) storage.CircularCategory {
	storageC := storage.CircularCategory{
		Name:        req.Name,
		Description: req.Description,
		Status:      int32(req.Status),
		Position:    int32(req.Position),
		CRUDTimeDate: storage.CRUDTimeDate{
			CreatedAt: req.CreatedAt.AsTime(),
			CreatedBy: "24182dde-5666-48f6-b38e-12f72477d9cc",
		},
	}

	return storageC
}
