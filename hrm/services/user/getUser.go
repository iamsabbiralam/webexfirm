package user

import (
	"context"

	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	user "practice/webex/gunk/v1/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Svc) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "Service.GetUser")
	res, err := h.core.GetUser(ctx, storage.User{
		Email: req.User.Email,
	})
	if err != nil {
		logging.WithError(err, log).Error("error with logging in user")
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return &user.GetUserResponse{
		User: &user.User{
			ID:       res.ID,
			Email:    res.Email,
			Password: res.Password,
		},
	}, nil
}
