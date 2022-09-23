package postgres

import (
	"context"
	"personal/webex/hrm/storage"
	"personal/webex/serviceutil/logging"
)

const getLogin = `
	SELECT
		email,
		password
	FROM
		users
	WHERE
		email = $1
	And 	deleted_at is null
	AND 	status = 1
`

func (s *Storage) Login(ctx context.Context, email string) (storage.SignUP, error) {
	log := logging.FromContext(ctx)
	var usr storage.SignUP
	err := s.db.Get(&usr, getLogin, email)
	if err != nil {
		logging.WithError(err, log).Error("get login prepare failed")
		return usr, err
	}

	return usr, nil
}
