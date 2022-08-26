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

	SignUP struct {
		ID        string    `db:"id"`
		FirstName string    `db:"first_name"`
		LastName  string    `db:"last_name"`
		Username  string    `db:"username"`
		Email     string    `db:"email"`
		Image     string    `db:"image"`
		Phone     string    `db:"phone"`
		Password  string    `db:"password"`
		Gender    int       `db:"gender"`
		DOB       time.Time `db:"dob"`
		Status    int       `db:"status"`
		CRUDTimeDate
	}
)
