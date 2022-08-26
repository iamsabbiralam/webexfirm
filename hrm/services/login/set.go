package login

import (
	"context"

	"personal/webex/hrm/storage"
	"personal/webex/serviceutil/logging"

	login "personal/webex/gunk/v1/login"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) Login(ctx context.Context, req *login.LoginRequest) (*login.LoginResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "Service.Login")
	if err := h.ls.Login(ctx, storage.SignUP{
		Email:    req.Login.Email,
		Password: req.Login.Password,
	}); err != nil {
		logging.WithError(err, log).Error("error with logging in user")
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return &login.LoginResponse{}, nil
}
