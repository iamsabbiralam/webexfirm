package jobType

import (
	"context"

	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/status"
)

func (s *CoreSvc) ListJobTypes(ctx context.Context, listJob storage.JobTypes) ([]storage.JobTypes, error) {
	log := logging.FromContext(ctx).WithField("method", "core.jpb-types.ListJobTypes")
	jobTypes, err := s.store.ListJobTypes(ctx, listJob)
	if err != nil {
		errMsg := "failed to get list of job types"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	return jobTypes, nil
}
