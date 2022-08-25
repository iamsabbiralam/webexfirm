package postgres

import (
	"context"
	"personal/webex/hrm/storage"
)

const register = `
	INSERT INTO users (
		first_name,
		last_name,
		username,
		email,
		image,
		phone,
		password,
		gender,
		dob,
		status
	) VALUES (
		:first_name,
		:last_name,
		:username,
		:email,
		:image,
		:phone,
		:password,
		:gender,
		:dob,
		:status
	) RETURNING
		id`

func(s *Storage) SignUP(ctx context.Context, user storage.SignUP) (string, error) {
	stmt, err := s.db.PrepareNamed(register)
	if err != nil {
		return "", err
	}

	var id string
	if err := stmt.Get(&id, user); err != nil {
		return "", err
	}

	return id, nil
}
