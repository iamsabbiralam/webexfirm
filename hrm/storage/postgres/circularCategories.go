package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"practice/webex/hrm/storage"
	"practice/webex/serviceutil/logging"

	"github.com/jmoiron/sqlx"
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
		id`

func (s *Storage) CreateCircularCategory(ctx context.Context, req storage.CircularCategory) (string, error) {
	log := logging.FromContext(ctx).WithField("method", "storage.circular-category.CreateCircularCategory")
	stmt, err := s.db.PrepareNamed(insertCircularCategory)
	if err != nil {
		errMsg := "failed to prepared insert circular category statement"
		log.WithError(err).Error(errMsg)
		return "", status.Error(status.Convert(err).Code(), errMsg)
	}

	var id string
	if err := stmt.Get(&id, req); err != nil {
		errMsg := "failed to execute insert circular category statement"
		log.WithError(err).Error(errMsg)
		return "", status.Error(status.Convert(err).Code(), errMsg)
	}

	return id, nil
}

const listCircularCategories = `
	WITH cnt AS (select count(*) as count FROM circular_categories WHERE deleted_at IS NULL)
	SELECT
		*,
		cnt.count
	FROM
		circular_categories
	LEFT JOIN
		cnt on true
	WHERE
		deleted_at IS NULL`

func (s *Storage) ListCircularCategory(ctx context.Context, cc storage.CircularCategory) ([]storage.CircularCategory, error) {
	log := logging.FromContext(ctx).WithField("method", "storage.circular-category.ListCircularCategories")
	searchQ := listCircularCategories
	inp := []interface{}{}
	if cc.SearchTerm != "" {
		searchQL := []string{}
		searchQL = append(searchQL, " AND (name ILIKE ? ) ")
		nm := fmt.Sprintf("%%%s%%", cc.SearchTerm)
		inp = append(inp, nm)
		searchQ += strings.Join(searchQL, " ")
	}

	if cc.Status != 0 {
		searchQ += " AND status = ?"
		inp = append(inp, cc.Status)
	}

	if cc.SortBy != "ASC" {
		cc.SortBy = "DESC"
		searchQ += " ORDER BY created_at " + cc.SortBy
	}

	if cc.Limit > 0 {
		searchQ += " LIMIT ?"
		inp = append(inp, cc.Limit)
	}

	if cc.Offset > 0 {
		searchQ += " OFFSET ?"
		inp = append(inp, cc.Offset)
	}

	fullQuery, args, err := sqlx.In(searchQ, inp...)
	if err != nil {
		errMsg := "failed to execute full query circular category"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	var circularCategory []storage.CircularCategory
	if err := s.db.Select(&circularCategory, s.db.Rebind(fullQuery), args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}

		errMsg := "failed to bind full query circular category"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	return circularCategory, nil
}

const getCircularCategory = `
	SELECT 
		*
	FROM
		circular_categories
	WHERE
		id = :id`

func (s *Storage) GetCircularCategory(ctx context.Context, id string) (*storage.CircularCategory, error) {
	log := logging.FromContext(ctx).WithField("method", "storage.circular-category.GetCircularCategory")
	stmt, err := s.db.PrepareNamed(getCircularCategory)
	if err != nil {
		errMsg := "failed to get circular category by id"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	var cc storage.CircularCategory
	cc.ID = id
	if err := stmt.Get(&cc, cc); err != nil {
		errMsg := "failed to execute circular category by id"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
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
		*`

func (s *Storage) UpdateCircularCategory(ctx context.Context, ctg storage.CircularCategory) (*storage.CircularCategory, error) {
	log := logging.FromContext(ctx).WithField("method", "storage.circular-category.UpdateCircularCategory")
	stmt, err := s.db.PrepareNamed(updateCircularCategory)
	if err != nil {
		errMsg := "failed to prepared update circular category statement"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	defer stmt.Close()
	var cc storage.CircularCategory
	if err := stmt.Get(&cc, ctg); err != nil {
		errMsg := "failed to execute update circular category by id"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	return &cc, nil
}

const deleteCircularCategory = `
	UPDATE
		circular_categories
	SET
		deleted_at = now(),
		deleted_by = :deleted_by
	WHERE 
		id = :id
	RETURNING 
		*`

func (s *Storage) DeleteCircularCategory(ctx context.Context, cc storage.CircularCategory) error {
	log := logging.FromContext(ctx).WithField("method", "storage.circular-category.DeleteCircularCategory")
	stmt, err := s.db.PrepareNamedContext(ctx, deleteCircularCategory)
	if err != nil {
		errMsg := "failed to prepared delete circular category"
		log.WithError(err).Error(errMsg)
		return status.Error(status.Convert(err).Code(), errMsg)
	}

	defer stmt.Close()
	arg := map[string]interface{}{
		"id":         cc.ID,
		"deleted_by": "b6ddbe32-3d7e-4828-b2d7-da9927846e6b",
	}
	if _, err := stmt.Exec(arg); err != nil {
		errMsg := "failed to execute delete circular category"
		log.WithError(err).Error(errMsg)
		return status.Error(status.Convert(err).Code(), errMsg)
	}

	return nil
}
