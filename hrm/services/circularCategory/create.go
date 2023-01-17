package circularCategory

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cc "practice/webex/gunk/v1/circularCategory"

	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"
)

func (h *Handler) CreateCircularCategory(ctx context.Context, req *cc.CreateCircularCategoryRequest) (*cc.CreateCircularCategoryResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.CircularCategory.create")
	dbPrm := storage.CircularCategory{
		Name:        req.Name,
		Description: req.Description,
		Status:      int32(req.Status),
		Position:    int32(req.Position),
		CRUDTimeDate: storage.CRUDTimeDate{
			CreatedAt: req.CreatedAt.AsTime(),
			CreatedBy: "24182dde-5666-48f6-b38e-12f72477d9cc",
		},
	}

	id, err := h.ccst.CreateCircularCategory(ctx, dbPrm)
	if err != nil {
		logging.WithError(err, log).Error("failed to create circular category")
		return nil, status.Error(codes.Internal, "failed to create circular category")
	}

	return &cc.CreateCircularCategoryResponse{
		ID: id,
	}, nil
}
