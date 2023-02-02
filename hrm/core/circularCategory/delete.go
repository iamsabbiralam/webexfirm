package circularCategory

import (
	"context"
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/status"
)

func (cs *CoreSvc) DeleteCircularCategory(ctx context.Context, req storage.CircularCategory) error {
	log := logging.FromContext(ctx).WithField("method", "core.circular-category.deleteCircularCategory")
	if err := cs.store.DeleteCircularCategory(ctx, req); err != nil {
		errMsg := "Circular Category deleted failed"
		log.WithError(err).Error(errMsg)
		return status.Error(status.Convert(err).Code(), errMsg)
	}
	
	return nil
}
