package signUp

import (
	"context"

	"personal/webex/hrm/storage"
	"personal/webex/serviceutil/logging"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) Register(ctx context.Context, signUp storage.SignUP) (string, error) {
	log := logging.FromContext(ctx).WithField("method", "Register")
	id, err := s.st.SignUP(ctx, signUp)
	if err != nil {
		logging.WithError(err, log).Error("error with registering user")
		return "", status.Error(codes.Internal, "processing failed")
	}

	return id, nil
}
