package user

import (
	"context"

	userG "practice/webex/gunk/v1/user"
	"practice/webex/hrm/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) CreateUser(ctx context.Context, req *userG.CreateUserRequest) (*userG.CreateUserResponse, error) {
	user := storage.User{
		FirstName: req.User.FirstName,
		LastName:  req.User.LastName,
		Email:     req.User.Email,
		Password:  req.User.Password,
		Status:    int(req.User.Status),
		CRUDTimeDate: storage.CRUDTimeDate{
			CreatedBy: timestamppb.Now().String(),
		},
	}
	id, err := s.core.CreateUser(context.Background(), user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	return &userG.CreateUserResponse{
		ID: id,
	}, nil
}
