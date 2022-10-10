package storage

type User struct {
	ID          string  `db:"id"`
	FullName    string `db:"full_name"`
	Email       string `db:"email"`
}
