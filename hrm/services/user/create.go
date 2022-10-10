package user

import (
	"context"
	"log"

	user "practice/webex/gunk/v1/user"
	"practice/webex/hrm/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) Create(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	log.Printf("Request Todo: %#v\n", req.GetUser())
	// Need to Validate request
	userC := storage.User{
		FullName: req.GetUser().FullName,
		Email:    req.GetUser().Email,
	}
	id, err := s.core.Create(context.Background(), userC)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}
	return &user.CreateUserResponse{
		ID: id,
	}, nil
}
