package circularCategory

import (
	"context"
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (cs *CoreSvc) ListCircularCategory(ctx context.Context, req storage.CircularCategory) ([]storage.CircularCategory, error) {
	log := logging.FromContext(ctx).WithField("method", "ListCircularCategory")
	res, err := cs.store.ListCircularCategory(ctx, req)
	if err != nil {
		logging.WithError(err, log).Error("Failed to list circular category")
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return res, nil
}
