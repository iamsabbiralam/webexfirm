package postgres

import (
	"context"
	"fmt"
	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const insertCircularCategory = `
INSERT INTO circular_categories (
	name, 
	description, 
	status,
	position,
	created_at,
	created_by
) VALUES (
	:name, 
	:description, 
	:status,
	:position,
	now(),
	:created_by
) RETURNING
	id
`

func (s *Storage) CreateCircularCategory(ctx context.Context, req storage.CircularCategory) (string, error) {
	stmt, err := s.db.PrepareNamed(insertCircularCategory)
	if err != nil {
		return "", err
	}

	var id string
	if err := stmt.Get(&id, req); err != nil {
		return "", err
	}

	return id, nil
}

func (s *Storage) ListCircularCategory(ctx context.Context, cc storage.CircularCategory) ([]storage.CircularCategory, error) {
	var sCC []storage.CircularCategory
	var status, search string
	if cc.SearchTerm != "" {
		search = fmt.Sprintf("(AND name ILIKE '%%' || '%s' || '%%')", cc.SearchTerm)
	}

	limit := ""
	if cc.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT NULLIF(%d, 0) OFFSET %d;", cc.Limit, cc.Offset)
	}

	if cc.Status == 1 {
		status = "AND status = 1"
	}

	listCircularCategories := fmt.Sprintf("WITH cnt AS (select count(*) as count FROM circular_categories WHERE deleted_at IS NULL %s %s) SELECT *, cnt.count FROM circular_categories left join cnt on true WHERE deleted_at IS NULL %s %s ORDER BY created_at DESC", search, status, search, status)
	fullQuery := listCircularCategories + limit
	if err := s.db.Select(&sCC, fullQuery); err != nil {
		return nil, fmt.Errorf("executing circular categories list: %s", err.Error())
	}

	return sCC, nil
}

const getCircularCategory = `
SELECT 
	*
FROM
	circular_categories
WHERE
	id = :id
`

func (s *Storage) GetCircularCategory(ctx context.Context, id string) (*storage.CircularCategory, error) {
	stmt, err := s.db.PrepareNamed(getCircularCategory)
	if err != nil {
		return &storage.CircularCategory{}, err
	}

	var cc storage.CircularCategory
	cc.ID = id
	if err := stmt.Get(&cc, cc); err != nil {
		return &storage.CircularCategory{}, err
	}

	return &cc, nil
}

const updateCircularCategory = `
UPDATE
	 circular_categories
SET
	name = :name,
    	description = :description,
    	status = :status,
	position = :position,
	updated_at = now(),
	updated_by = :updated_by
WHERE 
	id = :id
RETURNING
	*
`

func (s *Storage) UpdateCircularCategory(ctx context.Context, ctg storage.CircularCategory) (*storage.CircularCategory, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateCircularCategory)
	if err != nil {
		return &storage.CircularCategory{}, err
	}

	defer stmt.Close()
	var cc storage.CircularCategory
	if err := stmt.Get(&cc, ctg); err != nil {
		return &cc, nil
	}

	return &cc, nil
}

const deleteReturnCause = `
UPDATE
	circular_categories
SET
	deleted_at = now(),
	deleted_by = :deleted_by
WHERE 
	id = :id
RETURNING 
	*;
`

func (s *Storage) DeleteCircularCategory(ctx context.Context, scc storage.CircularCategory) error {
	log := logging.FromContext(ctx)
	stmt, err := s.db.PrepareNamedContext(ctx, deleteReturnCause)
	if err != nil {
		logging.WithError(err, log).Error("unable to delete circular category")
		return err
	}

	defer stmt.Close()
	var circular storage.CircularCategory
	if err := stmt.Get(&circular, scc); err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	return nil
}
