package storage

import (
	"database/sql"
	"time"
)

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
		ID        string `db:"id"`
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		Email     string `db:"email"`
		Password  string `db:"password"`
		Status    int    `db:"status"`
		CRUDTimeDate
	}
)
