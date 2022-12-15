package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	usr "practice/webex/gunk/v1/user"
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Svc) GetAllUsers(ctx context.Context, req *usr.GetAllUserRequest) (*usr.GetAllUserResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.User.GetAllUsers")

	userList, err := h.core.GetAllUsers(ctx, storage.User{
		SearchTerm: req.SearchTerm,
		Offset:     req.Offset,
		Limit:      req.Limit,
	})
	if err != nil {
		logging.WithError(err, log).Error("no users found")
		return nil, status.Error(codes.NotFound, "no users found")
	}

	list := make([]*usr.User, len(userList))
	for i, u := range userList {
		list[i] = &usr.User{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Status:    usr.Status(u.Status),
			CreatedAt: tspb.New(u.CreatedAt),
			CreatedBy: u.CreatedBy,
			UpdatedAt: tspb.New(u.UpdatedAt),
			UpdatedBy: u.UpdatedBy,
		}
	}
	var total int32
	if len(userList) > 0 {
		total = int32(userList[0].Count)
	}

	return &usr.GetAllUserResponse{
		User:  list,
		Total: total,
	}, nil
}
