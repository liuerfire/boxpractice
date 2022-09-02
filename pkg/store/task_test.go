package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/liuerfire/boxpractice/pkg/dto"
	"github.com/liuerfire/boxpractice/pkg/models"
)

func TestTask(t *testing.T) {
	store, cleanup := helperConnect(t)
	defer cleanup()

	ctx := context.Background()

	hospital, err := store.CreateHospital(ctx, &dto.Hospital{
		Name:        "task_hospital",
		DisplayName: "task hospital",
	})
	assert.NoError(t, err)

	employeeA, err := store.CreateEmployee(ctx, &dto.Employee{
		HospitalID: hospital.ID,
		Username:   "a",
	})
	assert.NoError(t, err)

	employeeB, err := store.CreateEmployee(ctx, &dto.Employee{
		HospitalID: hospital.ID,
		Username:   "b",
	})

	assert.NoError(t, err)

	t.Run("CRUD Operations", func(t *testing.T) {
		cases := []struct {
			Employee *models.Employee
			Tasks    []dto.Task
		}{
			{
				Employee: employeeA,
				Tasks: []dto.Task{
					{
						HospitalID:  hospital.ID,
						OwnerID:     employeeA.ID,
						Title:       "task a1",
						Description: "task a1 desc",
						Priority:    "LOW",
						Status:      "OPEN",
					},
					{
						HospitalID:  hospital.ID,
						OwnerID:     employeeA.ID,
						Title:       "task a2",
						Description: "task a2 desc",
						Priority:    "URGENT",
						Status:      "OPEN",
					},
				},
			},
			{
				Employee: employeeB,
				Tasks: []dto.Task{
					{
						HospitalID:  hospital.ID,
						OwnerID:     employeeB.ID,
						Title:       "task b1",
						Description: "task b1 desc",
						Priority:    "LOW",
						Status:      "OPEN",
					},
					{
						HospitalID:  hospital.ID,
						OwnerID:     employeeB.ID,
						Title:       "task b2",
						Description: "task b2 desc",
						Priority:    "URGENT",
						Status:      "OPEN",
					},
					{
						HospitalID:  hospital.ID,
						OwnerID:     employeeB.ID,
						Title:       "task b3",
						Description: "task b3 desc",
						Priority:    "HIGHT",
						Status:      "OPEN",
					},
				},
			},
		}
		employeeATasks := len(cases[0].Tasks)
		employeeBTasks := len(cases[1].Tasks)
		hospitalTasks := employeeATasks + employeeBTasks
		for _, cc := range cases {
			employee := cc.Employee
			for _, task := range cc.Tasks {
				d := task
				ta, err := store.CreateTask(ctx, &d)
				assert.NoError(t, err)
				assert.Greater(t, ta.ID, int64(0))
				assert.Equal(t, hospital.ID, ta.HospitalID)
				assert.Equal(t, employee.ID, ta.OwnerID)
				assert.Equal(t, task.Title, ta.Title)
				assert.Equal(t, task.Description, ta.Description)
				assert.Equal(t, task.Priority, ta.Priority)
				assert.Equal(t, task.Status, ta.Status)

				taskSame, err := store.GetTask(ctx, ta.ID)
				assert.NoError(t, err)
				assert.Equal(t, ta.ID, taskSame.ID)
				assert.Equal(t, ta.HospitalID, taskSame.HospitalID)
				assert.Equal(t, ta.OwnerID, taskSame.OwnerID)
				assert.Equal(t, ta.Title, taskSame.Title)
				assert.Equal(t, ta.Description, taskSame.Description)
				assert.Equal(t, ta.Priority, taskSame.Priority)
				assert.Equal(t, ta.Status, taskSame.Status)

				n, err := store.UpdateTask(ctx, &dto.Task{
					ID:          ta.ID,
					HospitalID:  ta.HospitalID,
					OwnerID:     ta.OwnerID,
					Title:       ta.Title,
					Description: ta.Description,
					Priority:    "LOW",
					Status:      "FAILED",
				})
				assert.NoError(t, err)
				assert.Equal(t, int64(1), n)

				taskNew, err := store.GetTask(ctx, ta.ID)
				assert.NoError(t, err)
				assert.Equal(t, ta.ID, taskNew.ID)
				assert.Equal(t, ta.HospitalID, taskNew.HospitalID)
				assert.Equal(t, ta.OwnerID, taskNew.OwnerID)
				assert.Equal(t, ta.Title, taskNew.Title)
				assert.Equal(t, ta.Description, taskNew.Description)
				assert.Equal(t, "LOW", taskNew.Priority)
				assert.Equal(t, "FAILED", taskNew.Status)
			}
		}
		tasks, err := store.FindTasksByHospital(ctx, hospital.ID, 0, 10)
		assert.NoError(t, err)
		assert.Equal(t, hospitalTasks, len(tasks))

		hn, err := store.CountTasksByHospital(ctx, hospital.ID)
		assert.NoError(t, err)
		assert.Equal(t, uint(hospitalTasks), hn)

		atasks, err := store.FindTasksByOwner(ctx, employeeA.ID, 0, 10)
		assert.NoError(t, err)
		assert.Equal(t, employeeATasks, len(atasks))

		an, err := store.CountTasksByOwner(ctx, employeeA.ID)
		assert.NoError(t, err)
		assert.Equal(t, uint(employeeATasks), an)

		btasks, err := store.FindTasksByOwner(ctx, employeeB.ID, 0, 10)
		assert.NoError(t, err)
		assert.Equal(t, employeeBTasks, len(btasks))

		bn, err := store.CountTasksByOwner(ctx, employeeB.ID)
		assert.NoError(t, err)
		assert.Equal(t, uint(employeeBTasks), bn)
	})
}
