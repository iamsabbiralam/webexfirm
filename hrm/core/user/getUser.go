package user

import (
	"context"

	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"
)

func (s *CoreSvc) GetUser(ctx context.Context, login storage.User) (storage.User, error) {
	log := logging.FromContext(ctx).WithField("method", "Login")
	res, err := s.store.GetUser(ctx, login)
	if err != nil {
		log.WithError(err).Error("failed to login")
		return storage.User{}, err
	}

	return res, nil
}
