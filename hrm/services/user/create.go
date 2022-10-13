package user

import (
	"context"
	"log"

	userG "practice/webex/gunk/v1/user"
	"practice/webex/hrm/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) CreateUser(ctx context.Context, req *userG.CreateUserRequest) (*userG.CreateUserResponse, error) {
	log.Printf("Request Todo: %#v\n", req.GetUser())
	// Need to Validate request
	user := storage.User{
		ID:           req.User.ID,
		FirstName:    req.User.FirstName,
		LastName:     req.User.LastName,
		Email:        req.User.Email,
		Password:     req.User.Password,
		UserName:     req.User.UserName,
		DOB:          req.User.DOB.String(),
		Gender:       req.User.Gender,
		PhoneNumber:  req.User.PhoneNumber,
		CRUDTimeDate: storage.CRUDTimeDate{
			CreatedBy: req.User.CreatedBy,
			UpdatedBy:  req.User.UpdatedBy,
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
