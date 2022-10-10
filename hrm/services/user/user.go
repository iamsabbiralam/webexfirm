package user

import (
	"context"

	"practice/webex/hrm/storage"
	user "practice/webex/gunk/v1/user"
)

type userCoreStore interface {
	Create(context.Context, storage.User) (string, error)
}

type Svc struct{
	user.UnimplementedUserServiceServer
	core userCoreStore
}

func NewUserServer(c userCoreStore) *Svc {
	return &Svc{
		core: c,
	}
}

