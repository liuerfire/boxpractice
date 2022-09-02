package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/liuerfire/boxpractice/pkg/dto"
	"github.com/liuerfire/boxpractice/pkg/models"
)

func TestEmployee(t *testing.T) {
	store, cleanup := helperConnect(t)
	defer cleanup()

	ctx := context.Background()

	hospital, err := store.CreateHospital(ctx, &dto.Hospital{
		Name:        "employ_hospital",
		DisplayName: "employ hospital",
	})
	assert.NoError(t, err)

	var employee, employeeSame, employeeOther *models.Employee

	t.Run("CreateEmployee", func(t *testing.T) {
		employee, err = store.CreateEmployee(ctx, &dto.Employee{
			HospitalID: hospital.ID,
			Username:   "alice",
			FirstName:  "Alice",
			LastName:   "Wong",
		})
		assert.NoError(t, err)
		assert.Greater(t, employee.ID, int64(0))
		assert.Equal(t, hospital.ID, employee.HospitalID)
		assert.Equal(t, "alice", employee.Username)
		assert.Equal(t, "Alice", employee.FirstName)
		assert.Equal(t, "Wong", employee.LastName)
	})

	t.Run("CreateAnotherEmployee", func(t *testing.T) {
		employeeOther, err = store.CreateEmployee(ctx, &dto.Employee{
			HospitalID: hospital.ID,
			Username:   "bob",
			FirstName:  "Bob",
			LastName:   "Hong",
		})
		assert.NoError(t, err)
		assert.Greater(t, employeeOther.ID, int64(0))
		assert.Equal(t, hospital.ID, employeeOther.HospitalID)
		assert.Equal(t, "bob", employeeOther.Username)
		assert.Equal(t, "Bob", employeeOther.FirstName)
		assert.Equal(t, "Hong", employeeOther.LastName)
	})

	t.Run("GetEmployee", func(t *testing.T) {
		employeeSame, err = store.GetEmployee(ctx, employee.ID)
		assert.NoError(t, err)
		assert.Equal(t, employee.ID, employeeSame.ID)
		assert.Equal(t, employee.HospitalID, employeeSame.HospitalID)
		assert.Equal(t, employee.Username, employeeSame.Username)
		assert.Equal(t, employee.FirstName, employeeSame.FirstName)
		assert.Equal(t, employee.LastName, employeeSame.LastName)
		assert.Equal(t, employee.CreatedAt, employee.CreatedAt)
		assert.Equal(t, employee.UpdatedAt, employee.UpdatedAt)
	})

	t.Run("FindEmployees", func(t *testing.T) {
		employees, err := store.FindEmployees(ctx, hospital.ID, 0, 10)
		assert.NoError(t, err)

		assert.Equal(t, 2, len(employees))

		assert.Equal(t, employee.ID, employees[0].ID)
		assert.Equal(t, employee.Username, employees[0].Username)
		assert.Equal(t, employee.FirstName, employees[0].FirstName)
		assert.Equal(t, employee.LastName, employees[0].LastName)

		assert.Equal(t, employeeOther.ID, employees[1].ID)
		assert.Equal(t, employeeOther.Username, employees[1].Username)
		assert.Equal(t, employeeOther.FirstName, employees[1].FirstName)
		assert.Equal(t, employeeOther.LastName, employees[1].LastName)
	})

	t.Run("ListEmployeesWithLimit", func(t *testing.T) {
		total, err := store.CountEmployees(ctx, hospital.ID)
		assert.NoError(t, err)
		assert.Equal(t, uint(2), total)

		employees, err := store.FindEmployees(ctx, hospital.ID, 0, 1)
		assert.NoError(t, err)

		assert.Equal(t, employee.ID, employees[0].ID)
		assert.Equal(t, employee.Username, employees[0].Username)
		assert.Equal(t, employee.FirstName, employees[0].FirstName)
		assert.Equal(t, employee.LastName, employees[0].LastName)

		employees, err = store.FindEmployees(ctx, hospital.ID, 1, 1)
		assert.NoError(t, err)

		assert.Equal(t, employeeOther.ID, employees[0].ID)
		assert.Equal(t, employeeOther.Username, employees[0].Username)
		assert.Equal(t, employeeOther.FirstName, employees[0].FirstName)
		assert.Equal(t, employeeOther.LastName, employees[0].LastName)
	})

	t.Run("CreateEmployeeIfExist", func(t *testing.T) {
		_, err = store.CreateEmployee(ctx, &dto.Employee{
			Username: employee.Username,
		})
		assert.True(t, IsErrDuplicateEntry(err))
	})
}
