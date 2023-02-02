package circularCategory

import (
	"context"

	"google.golang.org/grpc/status"

	cc "practice/webex/gunk/v1/circularCategory"
	"practice/webex/serviceutil/logging"
)

func (h *Handler) GetCircularCategory(ctx context.Context, req *cc.GetCircularCategoryRequest) (*cc.GetCircularCategoryResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.circular-category.GetCircularCategory")
	res, err := h.ccst.GetCircularCategory(ctx, req.GetID())
	if err != nil {
		errMsg := "no circular category found"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	return &cc.GetCircularCategoryResponse{
		ID:          res.ID,
		Name:        res.Name,
		Description: res.Description,
		Status:      cc.Status(res.Status),
		Position:    int64(res.Position),
	}, nil
}
