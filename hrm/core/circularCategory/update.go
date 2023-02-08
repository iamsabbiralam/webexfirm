package circularCategory

import (
	"context"

	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/status"
)

func (cs *CoreSvc) UpdateCircularCategory(ctx context.Context, req storage.CircularCategory) (*storage.CircularCategory, error) {
	log := logging.FromContext(ctx).WithField("method", "core.circular-category.UpdateCircularCategory")
	res, err := cs.store.UpdateCircularCategory(ctx, req)
	if err != nil {
		errMsg := "Failed to update circular category"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	return res, nil
}
