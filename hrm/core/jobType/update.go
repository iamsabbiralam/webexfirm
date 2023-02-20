package jobType

import (
	"context"
	
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/status"
)

func (s *CoreSvc) UpdateJobType(ctx context.Context, upJob storage.JobTypes) (string, error) {
	log := logging.FromContext(ctx).WithField("method", "core.job-type.UpdateJobType")
	id, err := s.store.UpdateJobType(ctx, upJob)
	if err != nil {
		errMsg := "Failed to update job type by id"
		log.WithError(err).Error(errMsg)
		return "", status.Error(status.Convert(err).Code(), errMsg)
	}

	return id, nil
}