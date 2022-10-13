package postgres

import (
	"context"
	"practice/webex/hrm/storage"
)

const insertUser = `
	INSERT INTO users(
		first_name,
		last_name,
		email,
		password,
		user_name,
		dob,
		gender,
		phone_number,
	) VALUES(
		:first_name,
		:last_name,
		:email,
		:password,
		:user_name,
		:dob,
		:gender,
		:phone_number,
	)
RETURNING
	id;
`
func (s *Storage) CreateUser(ctx context.Context, user storage.User) (string, error) {
	stmt, err := s.db.PrepareNamed(insertUser)
	if err != nil {
		return "", err
	}
	var id string
	if err := stmt.Get(&id, user); err != nil {
		return "", err
	}

	return id, nil
}
