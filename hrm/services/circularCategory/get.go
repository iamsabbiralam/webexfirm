package circularCategory

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cc "practice/webex/gunk/v1/circularCategory"
	"practice/webex/serviceutil/logging"
)

func (h *Handler) GetCircularCategory(ctx context.Context, req *cc.GetCircularCategoryRequest) (*cc.GetCircularCategoryResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.circularCategory.GetCircularCategory")
	log.Trace("request received")
	res, err := h.ccst.GetCircularCategory(ctx, req.GetID())
	if err != nil {
		logging.WithError(err, log).Error("no circular category found")
		return nil, status.Error(codes.NotFound, "circular category doesn't exists")
	}

	return &cc.GetCircularCategoryResponse{
		ID:          res.ID,
		Name:        res.Name,
		Description: res.Description,
		Status:      cc.Status(res.Status),
		Position:    int64(res.Position),
	}, nil
}
