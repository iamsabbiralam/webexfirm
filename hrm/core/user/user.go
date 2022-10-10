package user

import (
	"context"

	"practice/webex/hrm/storage"
)

type userStore interface {
	Create(context.Context, storage.User) (string, error)
}

type CoreSvc struct {
	store userStore
}

func NewCoreSvc(s userStore) *CoreSvc {
	return &CoreSvc{
		store: s,
	}
}

func (cs CoreSvc) Create(ctx context.Context, t storage.User) (string, error) {
	return cs.store.Create(ctx, t)
}