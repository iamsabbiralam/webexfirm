package circularCategory

import (
	"context"
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/status"
)

func (cs *CoreSvc) GetCircularCategory(ctx context.Context, id string) (*storage.CircularCategory, error) {
	log := logging.FromContext(ctx).WithField("method", "core.circular-category.GetCircularCategory")
	res, err := cs.store.GetCircularCategory(ctx, id)
	if err != nil && err != storage.NotFound {
		errMsg := "Failed to get Circular record"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	return res, nil
}
