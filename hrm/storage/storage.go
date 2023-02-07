package storage

import (
	"database/sql"
	"errors"
	"time"
)

// NotFound is returned when the requested resource does not exist.
var NotFound = errors.New("not found")

type (
	CRUDTimeDate struct {
		CreatedAt time.Time      `db:"created_at,omitempty"`
		CreatedBy string         `db:"created_by"`
		UpdatedAt time.Time      `db:"updated_at,omitempty"`
		UpdatedBy string         `db:"updated_by,omitempty"`
		DeletedAt sql.NullTime   `db:"deleted_at,omitempty"`
		DeletedBy sql.NullString `db:"deleted_by,omitempty"`
	}

	User struct {
		ID         string `db:"id"`
		FirstName  string `db:"first_name"`
		LastName   string `db:"last_name"`
		Email      string `db:"email"`
		Password   string `db:"password"`
		Status     int    `db:"status"`
		SearchTerm string `db:"search_term"`
		Offset     int32  `db:"offset"`
		Limit      int32  `db:"limit"`
		Count      int
		SortBy     string
		CRUDTimeDate
	}

	CircularCategory struct {
		ID          string `db:"id"`
		Name        string `db:"name"`
		Description string `db:"description"`
		Status      int32  `db:"status"`
		Position    int32  `db:"position,"`
		SearchTerm  string `db:"search_term"`
		Offset      int32  `db:"offset"`
		Limit       int32  `db:"limit"`
		Count       int
		SortBy      string
		CRUDTimeDate
	}

	JobTypes struct {
		ID         string `db:"id"`
		Name       string `db:"name"`
		Status     int    `db:"status"`
		Position   int    `db:"position"`
		SearchTerm string `db:"search_term"`
		Offset     int32  `db:"offset"`
		Limit      int32  `db:"limit"`
		Count      int
		CRUDTimeDate
	}
)
