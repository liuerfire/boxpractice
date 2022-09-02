package models

import (
	"time"
)

type Employee struct {
	ID         int64     `db:"id"`
	HospitalID int64     `db:"hospital_id"`
	Username   string    `db:"username"`
	FirstName  string    `db:"first_name"`
	LastName   string    `db:"last_name"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
