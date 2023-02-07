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

const insertJobTypes = `
	INSERT INTO job_types (
		name,
		status,
		position,
		created_at,
		created_by
	) VALUES (
		:name,
		:status,
		:position,
		now(),
		:created_by
	) RETURNING
		id`

func (s *Storage) CreateJobType(ctx context.Context, req storage.JobTypes) (string, error) {
	log := logging.FromContext(ctx).WithField("method", "storage.job-type.CreateJobType")
	stmt, err := s.db.PrepareNamed(insertJobTypes)
	if err != nil {
		errMsg := "failed to prepared insert job type statement"
		log.WithError(err).Error(errMsg)
		return "", status.Error(status.Convert(err).Code(), errMsg)
	}

	var id string
	if err := stmt.Get(&id, req); err != nil {
		errMsg := "failed to execute insert job type statement"
		log.WithError(err).Error(errMsg)
		return "", status.Error(status.Convert(err).Code(), errMsg)
	}

	return id, nil
}

const listJobTypes = `
	WITH cnt AS (select count(*) as count FROM job_types WHERE deleted_at IS NULL)
	SELECT
		*,
		cnt.count
	FROM
		job_types
	LEFT JOIN
		cnt on true
	WHERE
		deleted_at IS NULL`

func (s *Storage) ListJobTypes(ctx context.Context, job storage.JobTypes) ([]storage.JobTypes, error) {
	log := logging.FromContext(ctx).WithField("method", "storage.job-type.ListJobTypes")
	searchQ := listJobTypes
	inp := []interface{}{}
	if job.SearchTerm != "" {
		searchQL := []string{}
		searchQL = append(searchQL, " AND (name ILIKE ? ) ")
		nm := fmt.Sprintf("%%%s%%", job.SearchTerm)
		inp = append(inp, nm)
		searchQ += strings.Join(searchQL, " ")
	}

	if job.Status != 0 {
		searchQ += " AND status = ?"
		inp = append(inp, job.Status)
	}

	if job.Limit > 0 {
		searchQ += " LIMIT ?"
		inp = append(inp, job.Limit)
	}

	if job.Offset > 0 {
		searchQ += " OFFSET ?"
		inp = append(inp, job.Offset)
	}

	fullQuery, args, err := sqlx.In(searchQ, inp...)
	if err != nil {
		errMsg := "failed to execute full query job types"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	var jobTypes []storage.JobTypes
	if err := s.db.Select(&jobTypes, s.db.Rebind(fullQuery), args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}

		errMsg := "failed to bind full query job types"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	return jobTypes, nil
}

const getJobType = `
	SELECT 
		*
	FROM
		job_types
	WHERE
		id = :id`

func (s *Storage) GetJobType(ctx context.Context, id string) (*storage.JobTypes, error) {
	log := logging.FromContext(ctx).WithField("method", "storage.job-type.GetJobType")
	stmt, err := s.db.PrepareNamed(getJobType)
	if err != nil {
		errMsg := "failed to get job type by id"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	var job storage.JobTypes
	job.ID = id
	if err := stmt.Get(&job, job); err != nil {
		errMsg := "failed to execute job type by id"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	return &job, nil
}

const updateJobType = `
	UPDATE
		job_types
	SET
		name = COALESCE(NULLIF(:name, ''), name),
		status = COALESCE(NULLIF(:status, 0), status),
		position = COALESCE(NULLIF(:position, 0), position),
		updated_by = COALESCE(NULLIF(:updated_by, ''), updated_by),
		updated_at = now()
	WHERE 
		id = :id
	RETURNING
		*`

func (s *Storage) UpdateJobType(ctx context.Context, jobType storage.JobTypes) (*storage.JobTypes, error) {
	log := logging.FromContext(ctx).WithField("method", "storage.job-type.UpdateJobTypes")
	stmt, err := s.db.PrepareNamed(updateJobType)
	if err != nil {
		errMsg := "failed to prepared update job type statement"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}

	defer stmt.Close()
	var job storage.JobTypes
	if err := stmt.Get(&job, jobType); err != nil {
		errMsg := "failed to execute update job type by id"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(status.Convert(err).Code(), errMsg)
	}
	
	return &job, nil
}

const deleteJobType = `
	UPDATE
		job_types
	SET
		deleted_at = now(),
		deleted_by = :deleted_by
	WHERE 
		id = :id
	RETURNING 
		*`

func (s *Storage) DeleteJobType(ctx context.Context, cc storage.JobTypes) error {
	log := logging.FromContext(ctx).WithField("method", "storage.job-type.DeleteJobTypes")
	stmt, err := s.db.PrepareNamedContext(ctx, deleteJobType)
	if err != nil {
		errMsg := "failed to prepared delete job type"
		log.WithError(err).Error(errMsg)
		return status.Error(status.Convert(err).Code(), errMsg)
	}

	defer stmt.Close()
	arg := map[string]interface{}{
		"id": cc.ID,
		"deleted_by": "b6ddbe32-3d7e-4828-b2d7-da9927846e6b",
	}
	if _, err := stmt.Exec(arg); err != nil {
		errMsg := "failed to execute delete job type"
		log.WithError(err).Error(errMsg)
		return status.Error(status.Convert(err).Code(), errMsg)
	}

	return nil
}
