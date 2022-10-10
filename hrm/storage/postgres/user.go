package postgres

import (
	"context"
	"practice/webex/hrm/storage"
	"log"
)
const insertUser = `
	INSERT INTO users(
		full_name,
		email
	) VALUES(
		:full_name,
		:email
	)RETURNING id;
`
func (s *Storage) Create(ctx context.Context, t storage.User) (string, error) {
	stmt, err := s.db.PrepareNamed(insertUser)
	if err != nil {
		return "", err
	}
	var id string
	if err := stmt.Get(&id, t); err != nil {
		return "", err
	}
	log.Println("User ID: ", id)
	return id, nil

}	