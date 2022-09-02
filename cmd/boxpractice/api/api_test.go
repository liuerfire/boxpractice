package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/liuerfire/boxpractice/pkg/dto"
	"github.com/liuerfire/boxpractice/pkg/log"
	"github.com/liuerfire/boxpractice/pkg/models"
	"github.com/liuerfire/boxpractice/pkg/store"
)

func setup(t *testing.T) (api *API, cleanup func()) {
	t.Helper()
	assert := require.New(t)

	ctx := context.Background()
	logger := log.Init(10)

	sqlStore, err := store.NewSQLStore()
	assert.NoError(err)

	cleanup = func() {
		err := sqlStore.Cleanup()
		assert.NoError(err)
	}
	cleanup()

	api, err = InitAPIHandler(ctx, logger, sqlStore)
	assert.NoError(err)

	return
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// FIXME: refactor to use table-driven tests and add more cases.
func TestE2E(t *testing.T) {
	api, cleanup := setup(t)
	defer cleanup()

	router := mux.NewRouter()
	api.RegisterRouter(router)

	server := httptest.NewServer(router)
	defer server.Close()

	client := server.Client()

	hospital := dto.Hospital{
		Name:        "mine",
		DisplayName: "My Hospital",
	}

	t.Run("CreateHospital", func(t *testing.T) {
		data, _ := json.Marshal(hospital)

		resp, err := client.Post(fmt.Sprintf("%s/api/hospitals", server.URL), "application/json", bytes.NewReader(data))
		assert.NoError(t, err)
		defer resp.Body.Close()

		var tmp dto.Hospital
		err = json.NewDecoder(resp.Body).Decode(&tmp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, hospital.Name, tmp.Name)
		assert.Equal(t, hospital.DisplayName, tmp.DisplayName)
		assert.Greater(t, tmp.ID, int64(0))

		hospital = tmp
	})

	t.Run("GetHospital", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("%s/api/hospitals/%d", server.URL, hospital.ID))
		assert.NoError(t, err)
		defer resp.Body.Close()

		var h dto.Hospital
		err = json.NewDecoder(resp.Body).Decode(&h)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, hospital.Name, h.Name)
		assert.Equal(t, hospital.DisplayName, h.DisplayName)
		assert.Equal(t, hospital.ID, h.ID)
	})

	t.Run("ListHospitals", func(t *testing.T) {
		resp, err := client.Get(fmt.Sprintf("%s/api/hospitals", server.URL))
		assert.NoError(t, err)
		defer resp.Body.Close()

		var list dto.HospitalList
		err = json.NewDecoder(resp.Body).Decode(&list)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.GreaterOrEqual(t, list.Total, uint(0))
	})

	t.Run("UpdateHospital", func(t *testing.T) {
		path := fmt.Sprintf("%s/api/hospitals/%d", server.URL, hospital.ID)
		hospital.Name = "mine2"
		hospital.DisplayName = "blablabal"
		data, _ := json.Marshal(hospital)
		req, err := http.NewRequest("PUT", path, bytes.NewReader(data))
		assert.NoError(t, err)

		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		resp, err = client.Get(path)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var h dto.Hospital
		err = json.NewDecoder(resp.Body).Decode(&h)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, hospital.Name, h.Name)
		assert.Equal(t, hospital.DisplayName, h.DisplayName)
		assert.Equal(t, hospital.ID, h.ID)
	})

	employeeA := dto.Employee{
		HospitalID: hospital.ID,
		Username:   "aaaa",
		FirstName:  "aaaaa",
		LastName:   "bbbbb",
	}

	employeeB := dto.Employee{
		HospitalID: hospital.ID,
		Username:   "bbb",
		FirstName:  "bbbbb",
		LastName:   "dsdada",
	}

	t.Run("CreateEmployee", func(t *testing.T) {
		path := fmt.Sprintf("%s/api/hospitals/%d/employees", server.URL, hospital.ID)

		data, _ := json.Marshal(employeeA)
		resp, err := client.Post(path, "application/json", bytes.NewReader(data))
		assert.NoError(t, err)
		defer resp.Body.Close()

		var tmp dto.Employee

		err = json.NewDecoder(resp.Body).Decode(&tmp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, employeeA.HospitalID, tmp.HospitalID)
		assert.Equal(t, employeeA.Username, tmp.Username)
		assert.Equal(t, employeeA.FirstName, tmp.FirstName)
		assert.Equal(t, employeeA.LastName, tmp.LastName)

		employeeA = tmp

		data, _ = json.Marshal(employeeB)
		resp, err = client.Post(path, "application/json", bytes.NewReader(data))
		assert.NoError(t, err)
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&tmp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, employeeB.HospitalID, tmp.HospitalID)
		assert.Equal(t, employeeB.Username, tmp.Username)
		assert.Equal(t, employeeB.FirstName, tmp.FirstName)
		assert.Equal(t, employeeB.LastName, tmp.LastName)

		employeeB = tmp
	})

	t.Run("ListEmployees", func(t *testing.T) {
		path := fmt.Sprintf("%s/api/hospitals/%d/employees", server.URL, hospital.ID)

		resp, err := client.Get(path)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var list dto.EmployeeList
		err = json.NewDecoder(resp.Body).Decode(&list)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.GreaterOrEqual(t, list.Total, uint(2))
	})

	t.Run("GetEmployee", func(t *testing.T) {
		path := fmt.Sprintf("%s/api/employees/%d", server.URL, employeeA.ID)

		resp, err := client.Get(path)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var e dto.Employee
		err = json.NewDecoder(resp.Body).Decode(&e)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, employeeA.HospitalID, e.HospitalID)
		assert.Equal(t, employeeA.Username, e.Username)
		assert.Equal(t, employeeA.FirstName, e.FirstName)
		assert.Equal(t, employeeA.LastName, e.LastName)

		path = fmt.Sprintf("%s/api/employees/%d", server.URL, employeeB.ID)

		resp, err = client.Get(path)
		assert.NoError(t, err)
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&e)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, employeeB.HospitalID, e.HospitalID)
		assert.Equal(t, employeeB.Username, e.Username)
		assert.Equal(t, employeeB.FirstName, e.FirstName)
		assert.Equal(t, employeeB.LastName, e.LastName)
	})

	taskA := dto.Task{
		HospitalID:  hospital.ID,
		OwnerID:     employeeA.ID,
		Title:       "a",
		Description: "a desc",
		Priority:    models.TaskPriorityUrgent,
		Status:      models.TaskStatusOpen,
	}

	taskB := dto.Task{
		HospitalID:  hospital.ID,
		OwnerID:     employeeB.ID,
		Title:       "b",
		Description: "b desc",
		Priority:    models.TaskPriorityLow,
		Status:      models.TaskStatusOpen,
	}

	t.Run("CreateTask", func(t *testing.T) {
		path := fmt.Sprintf("%s/api/hospitals/%d/tasks", server.URL, hospital.ID)

		data, _ := json.Marshal(taskA)
		resp, err := client.Post(path, "application/json", bytes.NewReader(data))
		assert.NoError(t, err)
		defer resp.Body.Close()

		var tmp dto.Task

		err = json.NewDecoder(resp.Body).Decode(&tmp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, taskA.HospitalID, tmp.HospitalID)
		assert.Equal(t, taskA.OwnerID, tmp.OwnerID)
		assert.Equal(t, taskA.Title, tmp.Title)
		assert.Equal(t, taskA.Description, tmp.Description)
		assert.Equal(t, taskA.Priority, tmp.Priority)
		assert.Equal(t, taskA.Status, tmp.Status)

		taskA = tmp

		data, _ = json.Marshal(taskB)
		resp, err = client.Post(path, "application/json", bytes.NewReader(data))
		assert.NoError(t, err)
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&tmp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, taskB.HospitalID, tmp.HospitalID)
		assert.Equal(t, taskB.OwnerID, tmp.OwnerID)
		assert.Equal(t, taskB.Title, tmp.Title)
		assert.Equal(t, taskB.Description, tmp.Description)
		assert.Equal(t, taskB.Priority, tmp.Priority)
		assert.Equal(t, taskB.Status, tmp.Status)

		taskB = tmp
	})

	t.Run("ListTasks", func(t *testing.T) {
		path := fmt.Sprintf("%s/api/hospitals/%d/tasks", server.URL, hospital.ID)

		resp, err := client.Get(path)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var tmp dto.TaskList

		err = json.NewDecoder(resp.Body).Decode(&tmp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, uint(2), tmp.Total)

		path = fmt.Sprintf("%s/api/employees/%d/tasks", server.URL, employeeA.ID)

		resp, err = client.Get(path)
		assert.NoError(t, err)
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&tmp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, uint(1), tmp.Total)

		path = fmt.Sprintf("%s/api/employees/%d/tasks", server.URL, employeeB.ID)

		resp, err = client.Get(path)
		assert.NoError(t, err)
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&tmp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, uint(1), tmp.Total)
	})

	t.Run("AssignTask", func(t *testing.T) {
		path := fmt.Sprintf("%s/api/tasks/%d/assign", server.URL, taskA.ID)
		data, _ := json.Marshal(map[string]int64{"ownerId": taskB.ID})

		resp, err := client.Post(path, "application/json", bytes.NewReader(data))
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		path = fmt.Sprintf("%s/api/employees/%d/tasks", server.URL, employeeB.ID)

		resp, err = client.Get(path)
		assert.NoError(t, err)
		defer resp.Body.Close()

		var tmp dto.TaskList
		err = json.NewDecoder(resp.Body).Decode(&tmp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, uint(2), tmp.Total)
	})
}
