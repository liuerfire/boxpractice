package store

import (
	"context"
	"time"

	"github.com/liuerfire/boxpractice/pkg/dto"
	"github.com/liuerfire/boxpractice/pkg/models"
)

func (s *SQLStore) GetEmployee(ctx context.Context, id int64) (*models.Employee, error) {
	var e models.Employee
	sql := "select id, hospital_id, username, first_name, last_name, created_at, updated_at from employee where id = ?"
	err := s.db.GetContext(ctx, &e, sql, id)
	return &e, err
}

func (s *SQLStore) CreateEmployee(ctx context.Context, e *dto.Employee) (*models.Employee, error) {
	employee := &models.Employee{
		HospitalID: e.HospitalID,
		Username:   e.Username,
		FirstName:  e.FirstName,
		LastName:   e.LastName,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	sql := "insert into employee (hospital_id, username, first_name, last_name, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	r, err := s.db.ExecContext(ctx,
		sql, employee.HospitalID, employee.Username,
		employee.FirstName, employee.LastName,
		employee.CreatedAt, employee.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}
	employee.ID = id
	return employee, nil
}

func (s *SQLStore) FindEmployees(ctx context.Context, hid int64, offset, limit uint) ([]*models.Employee, error) {
	var employees []*models.Employee
	sql := "select id, username, first_name, last_name, created_at, updated_at from employee where hospital_id = ? order by id limit ?, ?"
	if err := s.db.Select(&employees, sql, hid, offset, limit); err != nil {
		return nil, err
	}
	return employees, nil
}

func (s *SQLStore) CountEmployees(ctx context.Context, hid int64) (uint, error) {
	var count uint
	sql := "select count(1) from employee where hospital_id = ?"
	if err := s.db.Get(&count, sql, hid); err != nil {
		return 0, err
	}
	return count, nil
}
