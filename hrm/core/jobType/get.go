package jobType

import (
	"context"

	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/status"
)

func (s *CoreSvc) GetJobType(ctx context.Context, id string) (*storage.JobTypes, error) {
	log := logging.FromContext(ctx).WithField("method", "core.job-type.GetJobTypes")
	getJob, err := s.store.GetJobType(ctx, id)
	if err != nil && err != storage.NotFound {
		errMsg := "Failed to get job type by id"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	return getJob, nil
}

/* func (s *CoreSvc) GetJobTypesByTitle(ctx context.Context, name string) (*storage.JobTypes, error) {
	log := logging.FromContext(ctx).WithField("method", "Core.GetJobTypesByTitle")
	m, err := s.jts.GetJobTypesByTitle(ctx, name)
	if err != nil {
		logging.WithError(err, log).Error("read storage")
		return nil, status.Error(codes.Internal, "processing failed")
	}
	return m, nil
} */