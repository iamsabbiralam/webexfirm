package circularCategory

import (
	"context"
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (cs *CoreSvc) GetCircularCategory(ctx context.Context, id string) (*storage.CircularCategory, error) {
	log := logging.FromContext(ctx).WithField("method", "GetCircularCategory")
	res, err := cs.store.GetCircularCategory(ctx, id)
	if err != nil && err != storage.NotFound {
		logging.WithError(err, log).Error("Failed to get Circular record")
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return res, nil
}
