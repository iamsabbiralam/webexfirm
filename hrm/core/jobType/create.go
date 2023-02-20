package jobType

import (
	"context"

	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/status"
)

func (s *CoreSvc) CreateJobType(ctx context.Context, job storage.JobTypes) (string, error) {
	log := logging.FromContext(ctx).WithField("method", "core.job-types.CreateJobTypes")
	id, err := s.store.CreateJobType(ctx, job)
	if err != nil {
		logging.WithError(err, log).Error("Failed to create JobType storage entry")
		errMsg := "Failed to create JobType storage entry"
		log.WithError(err).Error(errMsg)
		return "", status.Error(status.Convert(err).Code(), errMsg)
	}

	return id, nil
}
