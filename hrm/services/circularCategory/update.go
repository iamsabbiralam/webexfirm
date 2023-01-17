package circularCategory

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cc "practice/webex/gunk/v1/circularCategory"
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"
)

func (h *Handler) UpdateCircularCategory(ctx context.Context, req *cc.UpdateCircularCategoryRequest) (*cc.UpdateCircularCategoryResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.circularCategory.UpdateCircularCategory")
	res, err := h.ccst.UpdateCircularCategory(ctx, storage.CircularCategory{
		ID:          req.GetID(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Status:      int32(req.GetStatus()),
		Position:    int32(req.GetPosition()),
	})
	if err != nil {
		errMsg := "no circular category found"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(codes.NotFound, "circular category doesn't exists")
	}

	return &cc.UpdateCircularCategoryResponse{
		ID:          res.ID,
		Name:        res.Name,
		Description: res.Description,
		Status:      cc.Status(res.Status),
		Position:    int64(res.Position),
	}, nil
}
