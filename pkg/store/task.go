package store

import (
	"context"
	"time"

	"github.com/liuerfire/boxpractice/pkg/dto"
	"github.com/liuerfire/boxpractice/pkg/models"
)

func (s *SQLStore) GetTask(ctx context.Context, id int64) (*models.Task, error) {
	var t models.Task
	sql := "select id, hospital_id, owner_id, title, description, priority, status, created_at, updated_at from task where id = ?"
	err := s.db.GetContext(ctx, &t, sql, id)
	return &t, err
}

func (s *SQLStore) CreateTask(ctx context.Context, task *dto.Task) (*models.Task, error) {
	t := &models.Task{
		HospitalID:  task.HospitalID,
		OwnerID:     task.OwnerID,
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority,
		Status:      task.Status,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	sql := "insert into task (hospital_id, owner_id, title, description, priority, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	r, err := s.db.ExecContext(ctx, sql, t.HospitalID, t.OwnerID, t.Title, t.Description, t.Priority, t.Status, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}
	t.ID = id
	return t, nil
}

func (s *SQLStore) FindTasksByHospital(ctx context.Context, hosptialID int64, offset, limit uint) ([]*models.Task, error) {
	var tasks []*models.Task
	sql := "select id, hospital_id, owner_id, title, description, priority, status, created_at, updated_at from task where hospital_id = ? order by id limit ?, ?"
	if err := s.db.Select(&tasks, sql, hosptialID, offset, limit); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *SQLStore) CountTasksByHospital(ctx context.Context, hosptialID int64) (uint, error) {
	var count uint
	sql := "select count(1) from task where hospital_id = ?"
	if err := s.db.Get(&count, sql, hosptialID); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *SQLStore) FindTasksByOwner(ctx context.Context, oid int64, offset, limit uint) ([]*models.Task, error) {
	var tasks []*models.Task
	sql := "select id, title, description, priority, status, owner_id, created_at, updated_at from task where owner_id = ? order by id limit ?, ?"
	if err := s.db.Select(&tasks, sql, oid, offset, limit); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *SQLStore) CountTasksByOwner(ctx context.Context, oid int64) (uint, error) {
	var count uint
	sql := "select count(1) from task where owner_id = ?"
	if err := s.db.Get(&count, sql, oid); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *SQLStore) UpdateTask(ctx context.Context, task *dto.Task) (int64, error) {
	sql := "update task set owner_id=?, title=?, description=?, priority=?, status=? where id = ?"
	r, err := s.db.ExecContext(ctx, sql, task.OwnerID, task.Title, task.Description, task.Priority, task.Status, task.ID)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}
