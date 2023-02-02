package circularCategory

import (
	"context"
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/status"
)

func (cs *CoreSvc) CreateCircularCategory(ctx context.Context, req storage.CircularCategory) (string, error) {
	log := logging.FromContext(ctx).WithField("method", "core.circular-category.CreateCircularCategory")
	id, err := cs.store.CreateCircularCategory(ctx, req)
	if err != nil {
		errMsg := "Failed to create circular category storage entry"
		log.WithError(err).Error(errMsg)
		return "", status.Error(status.Convert(err).Code(), errMsg)
	}

	return id, nil
}
