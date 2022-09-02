package services

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	"github.com/liuerfire/boxpractice/pkg/dto"
	"github.com/liuerfire/boxpractice/pkg/store"
)

type EmployeeService struct {
	logger   logr.Logger
	sqlStore *store.SQLStore
}

func ProvideEmployeeService(logger logr.Logger, sqlStore *store.SQLStore) *EmployeeService {
	return &EmployeeService{
		logger:   logger.WithName("employeeService"),
		sqlStore: sqlStore,
	}
}

func (es *EmployeeService) CreateEmployee(ctx context.Context, e *dto.Employee) (*dto.Employee, error) {
	employee, err := es.sqlStore.CreateEmployee(ctx, e)
	if err != nil {
		if store.IsErrDuplicateEntry(err) {
			return nil, &ServiceError{ErrAlreadyExists, fmt.Sprintf("username exists: %s", e.Username)}
		}
		return nil, err
	}
	return &dto.Employee{
		ID:         employee.ID,
		HospitalID: employee.HospitalID,
		Username:   employee.Username,
		FirstName:  employee.FirstName,
		LastName:   employee.LastName,
		CreatedAt:  employee.CreatedAt,
	}, nil
}

func (es *EmployeeService) ListEmployees(ctx context.Context, id int64, page, limit uint) (*dto.EmployeeList, error) {
	total, err := es.sqlStore.CountEmployees(ctx, id)
	if err != nil {
		return nil, err
	}
	employees, err := es.sqlStore.FindEmployees(ctx, id, page, limit)
	if err != nil {
		return nil, err
	}
	items := make([]*dto.Employee, len(employees))
	for i := range employees {
		employee := employees[i]
		items[i] = &dto.Employee{
			ID:         employee.ID,
			HospitalID: employee.HospitalID,
			Username:   employee.Username,
			FirstName:  employee.FirstName,
			LastName:   employee.LastName,
			CreatedAt:  employee.CreatedAt,
		}
	}
	return &dto.EmployeeList{
		Total: total,
		Items: items,
	}, nil
}

func (es *EmployeeService) GetEmployee(ctx context.Context, id int64) (*dto.Employee, error) {
	employee, err := es.sqlStore.GetEmployee(ctx, id)
	if err != nil {
		if store.IsErrNotFound(err) {
			return nil, &ServiceError{ErrResourceNotFound, fmt.Sprintf("invalid id: %d", id)}
		}
		return nil, err
	}
	return &dto.Employee{
		ID:         employee.ID,
		HospitalID: employee.HospitalID,
		Username:   employee.Username,
		FirstName:  employee.FirstName,
		LastName:   employee.LastName,
		CreatedAt:  employee.CreatedAt,
	}, nil
}
