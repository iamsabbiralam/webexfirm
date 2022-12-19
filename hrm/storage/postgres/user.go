package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"practice/webex/hrm/storage"
	"strings"
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
		id,
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

func (s *Storage) GetAllUsers(ctx context.Context, user storage.User) ([]storage.User, error) {
	var users []storage.User
	var status string
	filterQL := []string{}
	filterQ := ""
	if user.SearchTerm != "" {
		srs := strings.Split(user.SearchTerm, " ")
		for _, v := range srs {
			filterQL = append(filterQL, fmt.Sprintf("(first_name ILIKE '%%' || '%s' || '%%' OR last_name ILIKE '%%' || '%s' || '%%' OR email ILIKE '%%' || '%s' || '%%')", v, v, v))
		}
	}

	if len(filterQL) > 0 {
		filterQ = "AND " + strings.Join(filterQL, "")
	}

	limit := ""
	if user.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT NULLIF(%d, 0) OFFSET %d;", user.Limit, user.Offset)
	}

	if user.Status != 0 {
		status = fmt.Sprintf("AND status = %d", user.Status)
	}

	listUsers := fmt.Sprintf(`WITH cnt AS (select count(*) as count FROM users WHERE deleted_at IS NULL %s) 
	SELECT 	
		id,
		first_name,
		last_name,
    		email,
		status,
		created_by,
		updated_by,
		cnt.count
	FROM users as u LEFT JOIN cnt on true WHERE deleted_at IS NULL %s ORDER BY created_at %s`, user.SearchTerm, status, filterQ)
	fullQuery := listUsers + limit
	if err := s.db.Select(&users, fullQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}

	return users, nil
}
