package dto

import (
	"time"
)

type Task struct {
	ID          int64     `json:"id,omitempty"`
	HospitalID  int64     `json:"HospitalId,omitempty"`
	OwnerID     int64     `json:"ownerId,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Priority    string    `json:"priority,omitempty"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}

type TaskList struct {
	Total uint    `json:"total"`
	Items []*Task `json:"items"`
}
