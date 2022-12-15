package user

import (
	"context"

	"practice/webex/hrm/storage"
	user "practice/webex/gunk/v1/user"
)

type userCoreStore interface {
	CreateUser(context.Context, storage.User) (string, error)
	GetUser(context.Context, storage.User) (storage.User, error)
	GetAllUsers(context.Context, storage.User) ([]storage.User, error)
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
