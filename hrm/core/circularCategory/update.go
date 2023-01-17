package circularCategory

import (
	"context"
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (cs *CoreSvc) UpdateCircularCategory(ctx context.Context, req storage.CircularCategory) (*storage.CircularCategory, error) {
	log := logging.FromContext(ctx).WithField("method", "UpdateCircularCategory")
	res, err := cs.store.UpdateCircularCategory(ctx, req)
	if err != nil {
		logging.WithError(err, log).Error("Failed to update circular category")
		return &storage.CircularCategory{}, status.Error(codes.Internal, "processing failed")
	}

	return res, nil
}
