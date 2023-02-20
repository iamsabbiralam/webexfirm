package jobType

import (
	"context"

	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/status"
)

func (s *CoreSvc) DeleteJobType(ctx context.Context, delJob storage.JobTypes) error {
	log := logging.FromContext(ctx).WithField("method", "core.job-type.DeleteJobType")
	if err := s.store.DeleteJobType(ctx, delJob); err != nil {
		errMsg := "Failed to delete jpb type"
		log.WithError(err).Error(errMsg)
		return status.Error(status.Convert(err).Code(), errMsg)
	}

	return nil
}
