package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"

	"github.com/liuerfire/boxpractice/internal/services"
)

type API struct {
	logger          logr.Logger
	hospitalService *services.HospitalService
	employeeService *services.EmployeeService
	taskService     *services.TaskService
}

func ProvideAPI(
	logger logr.Logger,
	hospitalService *services.HospitalService,
	employeeService *services.EmployeeService,
	taskService *services.TaskService,
) *API {
	return &API{
		logger:          logger.WithName("api"),
		hospitalService: hospitalService,
		employeeService: employeeService,
		taskService:     taskService,
	}
}

func (api *API) RegisterRouter(router *mux.Router) {
	r := router.PathPrefix("/api").Subrouter()
	r.Methods(http.MethodGet).Path("/hospitals").HandlerFunc(api.handleListHospitals)
	r.Methods(http.MethodPost).Path("/hospitals").HandlerFunc(api.handleCreateHospital)
	r.Methods(http.MethodGet).Path("/hospitals/{id}").HandlerFunc(api.handleGetHospital)
	r.Methods(http.MethodPut).Path("/hospitals/{id}").HandlerFunc(api.handleUpdateHospital)

	r.Methods(http.MethodGet).Path("/hospitals/{id}/employees").HandlerFunc(api.handleListEmployees)
	r.Methods(http.MethodPost).Path("/hospitals/{id}/employees").HandlerFunc(api.handleCreateEmployee)
	r.Methods(http.MethodGet).Path("/employees/{id}").HandlerFunc(api.handleGetEmployee)

	r.Methods(http.MethodGet).Path("/hospitals/{id}/tasks").HandlerFunc(api.handleListHospitalTasks)
	r.Methods(http.MethodGet).Path("/employees/{id}/tasks").HandlerFunc(api.handleListEmployeeTasks)
	r.Methods(http.MethodPost).Path("/hospitals/{id}/tasks").HandlerFunc(api.handleCreateTask)
	r.Methods(http.MethodPut).Path("/tasks/{id}").HandlerFunc(api.handleUpdateTask)
	r.Methods(http.MethodPost).Path("/tasks/{id}/assign").HandlerFunc(api.handleAssignTask)
}

func parsePaginationParams(pageStr, limitStr string) (uint, uint) {
	page := 0
	limit := 20
	p, _ := strconv.Atoi(pageStr)
	if p >= 1 {
		page = p
	}
	l, _ := strconv.Atoi(limitStr)
	if l > 0 {
		limit = l
	}
	if limit > 100 {
		limit = 100
	}
	return uint(page), uint(limit)
}

func renderJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		renderSvcError(w, err)
		return
	}
}

func renderBadRequestErr(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	rsp, _ := json.Marshal(newResponseError(services.ErrBadArgument, err.Error()))
	w.WriteHeader(http.StatusBadRequest)
	w.Write(rsp)
}

func renderSvcError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	var rsp []byte
	if svcErr, ok := err.(*services.ServiceError); ok {
		rsp, _ = json.Marshal(newResponseError(svcErr.Code, svcErr.Msg))
		w.WriteHeader(svcErr.Code.StatusCode())
		w.Write(rsp)
		return
	}
	rsp, _ = json.Marshal(newResponseError(services.ErrInternalError, err.Error()))
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(rsp)
}

func newResponseError(code services.ErrCode, msg string) map[string]any {
	return map[string]any{
		"error": map[string]any{
			"code": code,
			"msg":  msg,
		},
	}
}
