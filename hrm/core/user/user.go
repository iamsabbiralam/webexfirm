package user

import (
	"context"

	"practice/webex/hrm/storage"
)

type userStore interface {
	CreateUser(context.Context, storage.User) (string, error)
}

type CoreSvc struct {
	store userStore
}

func NewCoreSvc(s userStore) *CoreSvc {
	return &CoreSvc{
		store: s,
	}
}
