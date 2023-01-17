package circularCategory

import (
	"context"
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (cs *CoreSvc) DeleteCircularCategory(ctx context.Context, req storage.CircularCategory) error {
	log := logging.FromContext(ctx).WithField("method", "core.deleteCircularCategory")
	if err := cs.store.DeleteCircularCategory(ctx, req); err != nil {
		logging.WithError(err, log).Error("Circular Category deleted failed")
		return status.Error(codes.Internal, "processing failed")
	}
	
	return nil
}
