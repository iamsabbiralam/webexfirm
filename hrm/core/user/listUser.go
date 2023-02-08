package user

import (
	"context"

	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreSvc) GetAllUsers(ctx context.Context, user storage.User) ([]storage.User, error) {
	log := logging.FromContext(ctx).WithField("method", "GetAllUsers")
	res, err := s.store.GetAllUsers(ctx, user)
	if err != nil {
		logging.WithError(err, log).Error("Failed to list GetAllUsers")
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return res, nil
}
