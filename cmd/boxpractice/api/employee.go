package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/liuerfire/boxpractice/pkg/dto"
)

func (api *API) handleListEmployees(w http.ResponseWriter, r *http.Request) {
	hidStr := mux.Vars(r)["id"]
	hid, err := strconv.ParseInt(hidStr, 10, 64)
	if err != nil {
		renderBadRequestErr(w, err)
		return
	}
	page, limit := parsePaginationParams(r.URL.Query().Get("page"), r.URL.Query().Get("limit"))
	employeeList, err := api.employeeService.ListEmployees(r.Context(), hid, page-1, limit)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	renderJSON(w, http.StatusOK, employeeList)
}

func (api *API) handleCreateEmployee(w http.ResponseWriter, r *http.Request) {
	hidStr := mux.Vars(r)["id"]
	hid, err := strconv.ParseInt(hidStr, 10, 64)
	if err != nil {
		renderBadRequestErr(w, err)
		return
	}
	var req dto.Employee
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderBadRequestErr(w, err)
		return
	}
	if req.Username == "" {
		renderBadRequestErr(w, errors.New("username is null"))
		return
	}
	_, err = api.hospitalService.GetHospital(r.Context(), hid)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	req.HospitalID = hid
	employee, err := api.employeeService.CreateEmployee(r.Context(), &req)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	renderJSON(w, http.StatusCreated, employee)
}

func (api *API) handleGetEmployee(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		renderBadRequestErr(w, err)
		return
	}
	employee, err := api.employeeService.GetEmployee(r.Context(), id)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	renderJSON(w, http.StatusOK, employee)
}
