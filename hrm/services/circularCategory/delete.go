package circularCategory

import (
	"context"

	"google.golang.org/grpc/status"

	cc "practice/webex/gunk/v1/circularCategory"
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"
)

func (h *Handler) DeleteCircularCategory(ctx context.Context, req *cc.DeleteCircularCategoryRequest) (*cc.DeleteCircularCategoryResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.circular-category.DeleteCircularCategory")
	err := h.ccst.DeleteCircularCategory(ctx, storage.CircularCategory{
		ID: req.GetID(),
	})
	if err != nil {
		errMsg := "failed to delete circular category"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	return &cc.DeleteCircularCategoryResponse{}, nil
}
