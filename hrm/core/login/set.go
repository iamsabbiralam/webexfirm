package login

import (
	"context"
	"personal/webex/hrm/storage"
	"personal/webex/serviceutil/logging"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) Login(ctx context.Context, email string) (storage.SignUP, error) {
	log := logging.FromContext(ctx).WithField("method", "Login")
	res, err := s.st.Login(ctx, email)
	if err != nil {
		logging.WithError(err, log).Error("login record")
		return res, status.Error(codes.Internal, "processing failed")
	}

	return res, nil
}
