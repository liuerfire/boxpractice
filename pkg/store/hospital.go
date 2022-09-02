package store

import (
	"context"
	"time"

	"github.com/liuerfire/boxpractice/pkg/dto"
	"github.com/liuerfire/boxpractice/pkg/models"
)

func (s *SQLStore) GetHospital(ctx context.Context, id int64) (*models.Hospital, error) {
	var hospital models.Hospital
	sql := "select id, name, display_name, created_at, updated_at from hospital where id = ?"
	err := s.db.GetContext(ctx, &hospital, sql, id)
	return &hospital, err
}

func (s *SQLStore) CreateHospital(ctx context.Context, h *dto.Hospital) (*models.Hospital, error) {
	hs := &models.Hospital{
		Name:        h.Name,
		DisplayName: h.DisplayName,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	sql := "insert into hospital (name, display_name, created_at, updated_at) VALUES (?, ?, ?, ?)"
	r, err := s.db.ExecContext(ctx, sql, hs.Name, hs.DisplayName, hs.CreatedAt, hs.UpdatedAt)
	if err != nil {
		return nil, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return nil, err
	}
	hs.ID = id
	return hs, nil
}

func (s *SQLStore) UpdateHospital(ctx context.Context, h *dto.Hospital) (int64, error) {
	sql := "update hospital set name=?, display_name=? where id = ?"
	r, err := s.db.ExecContext(ctx, sql, h.Name, h.DisplayName, h.ID)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}

func (s *SQLStore) FindHospitals(ctx context.Context, offset, limit uint) ([]*models.Hospital, error) {
	var hospitals []*models.Hospital
	sql := "select id, name, display_name, created_at, updated_at from hospital order by id limit ?, ?"
	if err := s.db.Select(&hospitals, sql, offset, limit); err != nil {
		return nil, err
	}
	return hospitals, nil
}

func (s *SQLStore) CountHosptials(ctx context.Context) (uint, error) {
	var count uint
	sql := "select count(1) from hospital"
	if err := s.db.Get(&count, sql); err != nil {
		return 0, err
	}
	return count, nil
}
