package services

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	"github.com/liuerfire/boxpractice/pkg/dto"
	"github.com/liuerfire/boxpractice/pkg/store"
)

type TaskService struct {
	logger   logr.Logger
	sqlStore *store.SQLStore
}

func ProvideTaskService(logger logr.Logger, sqlStore *store.SQLStore) *TaskService {
	return &TaskService{
		logger:   logger.WithName("taskService"),
		sqlStore: sqlStore,
	}
}

func (ts *TaskService) CreateTask(ctx context.Context, t *dto.Task) (*dto.Task, error) {
	task, err := ts.sqlStore.CreateTask(ctx, t)
	if err != nil {
		return nil, err
	}
	return &dto.Task{
		ID:          task.ID,
		HospitalID:  task.HospitalID,
		OwnerID:     task.OwnerID,
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
	}, nil
}

func (ts *TaskService) ListTasksByHospital(ctx context.Context, hid int64, page, limit uint) (*dto.TaskList, error) {
	total, err := ts.sqlStore.CountTasksByHospital(ctx, hid)
	if err != nil {
		return nil, err
	}
	tasks, err := ts.sqlStore.FindTasksByHospital(ctx, hid, page, limit)
	if err != nil {
		return nil, err
	}
	items := make([]*dto.Task, len(tasks))
	for i := range tasks {
		employee := tasks[i]
		items[i] = &dto.Task{
			ID:          employee.ID,
			HospitalID:  employee.HospitalID,
			OwnerID:     employee.OwnerID,
			Title:       employee.Title,
			Description: employee.Description,
			Priority:    employee.Priority,
			Status:      employee.Status,
			CreatedAt:   employee.CreatedAt,
		}
	}
	return &dto.TaskList{
		Total: total,
		Items: items,
	}, nil
}

func (ts *TaskService) ListTasksByOwner(ctx context.Context, oid int64, page, limit uint) (*dto.TaskList, error) {
	total, err := ts.sqlStore.CountTasksByOwner(ctx, oid)
	if err != nil {
		return nil, err
	}
	tasks, err := ts.sqlStore.FindTasksByOwner(ctx, oid, page, limit)
	if err != nil {
		return nil, err
	}
	items := make([]*dto.Task, len(tasks))
	for i := range tasks {
		employee := tasks[i]
		items[i] = &dto.Task{
			ID:          employee.ID,
			HospitalID:  employee.HospitalID,
			OwnerID:     employee.OwnerID,
			Title:       employee.Title,
			Description: employee.Description,
			Priority:    employee.Priority,
			Status:      employee.Status,
			CreatedAt:   employee.CreatedAt,
		}
	}
	return &dto.TaskList{
		Total: total,
		Items: items,
	}, nil
}

func (ts *TaskService) GetTask(ctx context.Context, id int64) (*dto.Task, error) {
	task, err := ts.sqlStore.GetTask(ctx, id)
	if err != nil {
		if store.IsErrNotFound(err) {
			return nil, &ServiceError{ErrResourceNotFound, fmt.Sprintf("invalid id: %d", id)}
		}
		return nil, err
	}
	return &dto.Task{
		ID:          task.ID,
		HospitalID:  task.HospitalID,
		OwnerID:     task.OwnerID,
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
	}, nil
}

func (ts *TaskService) UpdateTask(ctx context.Context, t *dto.Task) error {
	r, err := ts.sqlStore.UpdateTask(ctx, t)
	if r == 0 {
		return &ServiceError{ErrResourceNotFound, fmt.Sprintf("invalid id: %d", t.ID)}
	}
	return err
}
