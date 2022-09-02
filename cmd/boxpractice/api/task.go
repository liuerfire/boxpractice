package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/liuerfire/boxpractice/internal/services"
	"github.com/liuerfire/boxpractice/pkg/dto"
	"github.com/liuerfire/boxpractice/pkg/models"
)

func (api *API) handleListHospitalTasks(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePaginationParams(r.URL.Query().Get("page"), r.URL.Query().Get("limit"))
	hidStr := mux.Vars(r)["id"]
	hid, err := strconv.ParseInt(hidStr, 10, 64)
	if err != nil {
		renderBadRequestErr(w, err)
		return
	}
	_, err = api.hospitalService.GetHospital(r.Context(), hid)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	taskList, err := api.taskService.ListTasksByHospital(r.Context(), hid, page-1, limit)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	renderJSON(w, http.StatusOK, taskList)
}

func (api *API) handleListEmployeeTasks(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePaginationParams(r.URL.Query().Get("page"), r.URL.Query().Get("limit"))
	idStr := mux.Vars(r)["id"]
	oid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		renderBadRequestErr(w, err)
		return
	}
	_, err = api.employeeService.GetEmployee(r.Context(), oid)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	taskList, err := api.taskService.ListTasksByOwner(r.Context(), oid, page-1, limit)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	renderJSON(w, http.StatusOK, taskList)
}

func (api *API) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	hidStr := mux.Vars(r)["id"]
	hid, err := strconv.ParseInt(hidStr, 10, 64)
	if err != nil {
		renderBadRequestErr(w, err)
		return
	}
	var req dto.Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderBadRequestErr(w, err)
		return
	}
	if err = validateTask(&req); err != nil {
		renderBadRequestErr(w, err)
		return
	}
	_, err = api.hospitalService.GetHospital(r.Context(), hid)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	owner, err := api.employeeService.GetEmployee(r.Context(), req.OwnerID)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	if owner.HospitalID != hid {
		renderSvcError(w, &services.ServiceError{Code: services.ErrPermissionDenied, Msg: "forbidden"})
		return
	}
	// The initial status
	req.Status = models.TaskStatusOpen
	task, err := api.taskService.CreateTask(r.Context(), &req)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	renderJSON(w, http.StatusCreated, task)
}

func (api *API) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		renderBadRequestErr(w, err)
		return
	}
	var req dto.Task
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderBadRequestErr(w, err)
		return
	}
	if err = validateTask(&req); err != nil {
		renderBadRequestErr(w, err)
		return
	}
	task, err := api.taskService.GetTask(r.Context(), id)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	task.Title = req.Title
	task.Description = req.Description
	task.Priority = req.Priority
	task.Status = req.Status
	if err := api.taskService.UpdateTask(r.Context(), task); err != nil {
		renderSvcError(w, err)
		return
	}
}

type assignTaskReq struct {
	OwnerID int64 `json:"ownerId"`
}

func (api *API) handleAssignTask(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		renderBadRequestErr(w, err)
		return
	}
	var req assignTaskReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderBadRequestErr(w, err)
		return
	}
	task, err := api.taskService.GetTask(r.Context(), id)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	employee, err := api.employeeService.GetEmployee(r.Context(), req.OwnerID)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	if task.HospitalID != employee.HospitalID {
		renderSvcError(w, &services.ServiceError{Code: services.ErrPermissionDenied, Msg: "forbidden"})
		return
	}
	task.OwnerID = employee.ID
	if err := api.taskService.UpdateTask(r.Context(), task); err != nil {
		renderSvcError(w, err)
		return
	}
}

func validateTask(t *dto.Task) error {
	if t.OwnerID <= 0 {
		return errors.New("invalid owner id")
	}
	if t.Title == "" {
		return errors.New("invalid title")
	}
	if !isValidPriority(t.Priority) {
		return errors.New("invalid priority")
	}
	if !isValidStatus(t.Status) {
		return errors.New("invalid status")
	}
	return nil
}

func isValidPriority(p string) bool {
	for _, elem := range []string{models.TaskPriorityUrgent, models.TaskPriorityHight, models.TaskPriorityLow} {
		if elem == p {
			return true
		}
	}
	return false
}

func isValidStatus(s string) bool {
	for _, elem := range []string{models.TaskStatusOpen, models.TaskStatusFAILED, models.TaskStatusCOMPLETED} {
		if elem == s {
			return true
		}
	}
	return false
}
