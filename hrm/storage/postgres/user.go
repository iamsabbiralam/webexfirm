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
		status
	) VALUES(
		:first_name,
		:last_name,
		:email,
		:password,
		:status
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

const getUser = `
	SELECT 
		email,
		password
	FROM
		users
	WHERE
		email = :email AND
		status = 1 AND
		deleted_at IS NULL
`
func (s *Storage) GetUser(ctx context.Context, user storage.User) (storage.User, error) {
	stmt, err := s.db.PrepareNamed(getUser)
	if err != nil {
		return storage.User{}, err
	}

	var getUser storage.User
	if err := stmt.Get(&getUser, user); err != nil {
		return storage.User{}, err
	}

	return getUser, nil
}
