package circularCategory

import (
	"context"
	
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/status"
)

func (cs *CoreSvc) ListCircularCategory(ctx context.Context, req storage.CircularCategory) ([]storage.CircularCategory, error) {
	log := logging.FromContext(ctx).WithField("method", "core.circular-category.ListCircularCategories")
	res, err := cs.store.ListCircularCategory(ctx, req)
	if err != nil {
		errMsg := "Failed to list circular category"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	return res, nil
}
