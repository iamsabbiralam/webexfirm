package login

import (
	"context"
	"personal/webex/hrm/storage"
	"personal/webex/serviceutil/logging"
)

func (s *Svc) Login(ctx context.Context, login storage.SignUP) error {
	log := logging.FromContext(ctx).WithField("method", "Login")
	if err := s.st.Login(ctx, login); err != nil {
		log.WithError(err).Error("failed to login")
		return err
	}

	return nil
}
