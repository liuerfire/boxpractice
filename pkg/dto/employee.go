package dto

import (
	"time"
)

type Employee struct {
	ID         int64     `json:"id,omitempty"`
	HospitalID int64     `json:"hospitalId,omitempty"`
	Username   string    `json:"username,omitempty"`
	FirstName  string    `json:"firstName,omitempty"`
	LastName   string    `json:"lastName,omitempty"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
}

type EmployeeList struct {
	Total uint        `json:"total"`
	Items []*Employee `json:"items"`
}
