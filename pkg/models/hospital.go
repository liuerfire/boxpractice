package models

import (
	"time"
)

type Hospital struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	DisplayName string    `db:"display_name"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
