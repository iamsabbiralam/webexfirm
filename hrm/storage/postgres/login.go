package postgres

import (
	"context"
	"personal/webex/hrm/storage"
)

const login = `
	SELECT
		email,
		password
	FROM
		users
	WHERE
		email = :email
	AND 	password = :password
	And 	deleted_at is null
	AND 	status = 1
`

func (s *Storage) Login(ctx context.Context, user storage.SignUP) error {
	stmt, err := s.db.PrepareNamed(login)
	if err != nil {
		return err
	}

	if err := stmt.Get(&user, user); err != nil {
		return err
	}

	return nil
}
