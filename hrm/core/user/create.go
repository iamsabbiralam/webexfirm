package user

import (
	"context"
	"practice/webex/hrm/storage"
)

func (cs *CoreSvc) CreateUser(ctx context.Context, user storage.User) (string, error) {
	id, err := cs.store.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	return id, nil
}
