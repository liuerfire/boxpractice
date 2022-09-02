package services

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	"github.com/liuerfire/boxpractice/pkg/dto"
	"github.com/liuerfire/boxpractice/pkg/store"
)

type HospitalService struct {
	logger   logr.Logger
	sqlStore *store.SQLStore
}

func ProvideHospitalService(logger logr.Logger, sqlStore *store.SQLStore) *HospitalService {
	return &HospitalService{
		logger:   logger.WithName("hospitalService"),
		sqlStore: sqlStore,
	}
}

func (hs *HospitalService) CreateHospital(ctx context.Context, h *dto.Hospital) (*dto.Hospital, error) {
	hospital, err := hs.sqlStore.CreateHospital(ctx, h)
	if err != nil {
		if store.IsErrDuplicateEntry(err) {
			return nil, &ServiceError{ErrAlreadyExists, fmt.Sprintf("username exists: %s", h.Name)}
		}
		return nil, err
	}
	return &dto.Hospital{
		ID:          hospital.ID,
		Name:        hospital.Name,
		DisplayName: hospital.DisplayName,
		CreatedAt:   hospital.CreatedAt,
	}, nil
}

func (hs *HospitalService) ListHospitals(ctx context.Context, page, limit uint) (*dto.HospitalList, error) {
	total, err := hs.sqlStore.CountHosptials(ctx)
	if err != nil {
		return nil, err
	}
	hospitals, err := hs.sqlStore.FindHospitals(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	items := make([]*dto.Hospital, len(hospitals))
	for i := range hospitals {
		hospital := hospitals[i]
		items[i] = &dto.Hospital{
			ID:          hospital.ID,
			Name:        hospital.Name,
			DisplayName: hospital.DisplayName,
			CreatedAt:   hospital.CreatedAt,
		}
	}
	return &dto.HospitalList{
		Total: total,
		Items: items,
	}, nil
}

func (hs *HospitalService) GetHospital(ctx context.Context, hid int64) (*dto.Hospital, error) {
	hospital, err := hs.sqlStore.GetHospital(ctx, hid)
	if err != nil {
		if store.IsErrNotFound(err) {
			return nil, &ServiceError{ErrResourceNotFound, fmt.Sprintf("invalid id: %d", hid)}
		}
		return nil, err
	}
	return &dto.Hospital{
		ID:          hospital.ID,
		Name:        hospital.Name,
		DisplayName: hospital.DisplayName,
		CreatedAt:   hospital.CreatedAt,
	}, nil
}

func (hs *HospitalService) UpdateHospital(ctx context.Context, h *dto.Hospital) error {
	r, err := hs.sqlStore.UpdateHospital(ctx, h)
	if r == 0 {
		return &ServiceError{ErrResourceNotFound, fmt.Sprintf("invalid id: %d", h.ID)}
	}
	return err
}
