package models

import (
	"time"
)

const (
	TaskPriorityUrgent = "URGENT"
	TaskPriorityHight  = "HIGHT"
	TaskPriorityLow    = "LOW"

	TaskStatusOpen      = "OPEN"
	TaskStatusFAILED    = "FAILED"
	TaskStatusCOMPLETED = "COMPLETED"
)

type Task struct {
	ID          int64     `db:"id"`
	HospitalID  int64     `db:"hospital_id"`
	OwnerID     int64     `db:"owner_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Priority    string    `db:"priority"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
