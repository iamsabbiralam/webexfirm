package signUp

import (
	"context"

	"personal/webex/hrm/storage"
	"personal/webex/serviceutil/logging"

	signup "personal/webex/gunk/v1/signUp"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) Registration(ctx context.Context, req *signup.RegisterRequest) (*signup.RegisterResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "Service.Registration")
	res, err := h.su.Register(ctx, storage.SignUP{
		FirstName: req.SignUP.FirstName,
		LastName:  req.SignUP.LastName,
		Username:  req.SignUP.Username,
		Email:     req.SignUP.Email,
		Image:     req.SignUP.Image,
		Phone:     req.SignUP.Phone,
		Password:  req.SignUP.Password,
		Gender:    int(req.SignUP.Gender),
		DOB:       req.SignUP.DOB.AsTime(),
		Status:    int(req.SignUP.Status),
	})
	if err != nil {
		logging.WithError(err, log).Error("error with registering user")
		return nil, status.Error(codes.Internal, "processing failed")
	}

	return &signup.RegisterResponse{
		ID: res,
	}, nil
}
