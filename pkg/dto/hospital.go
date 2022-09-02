package dto

import (
	"time"
)

type Hospital struct {
	ID          int64     `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	DisplayName string    `json:"displayName,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
}

type HospitalList struct {
	Total uint        `json:"total"`
	Items []*Hospital `json:"items"`
}
