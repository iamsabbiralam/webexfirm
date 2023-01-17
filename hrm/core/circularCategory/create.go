package circularCategory

import (
	"context"
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (cs *CoreSvc) CreateCircularCategory(ctx context.Context, req storage.CircularCategory) (string, error) {
	log := logging.FromContext(ctx).WithField("method", "CreateCircularCategory")

	id, err := cs.store.CreateCircularCategory(ctx, req)
	if err != nil {
		logging.WithError(err, log).Error("Failed to create circular category storage entry")
		return "", status.Error(codes.Internal, "processing failed")
	}

	return id, nil
}
