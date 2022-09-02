package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/liuerfire/boxpractice/pkg/dto"
)

func (api *API) handleListHospitals(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePaginationParams(r.URL.Query().Get("page"), r.URL.Query().Get("limit"))
	h, err := api.hospitalService.ListHospitals(r.Context(), page-1, limit)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	renderJSON(w, http.StatusOK, h)
}

func (api *API) handleCreateHospital(w http.ResponseWriter, r *http.Request) {
	var req dto.Hospital
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderBadRequestErr(w, err)
		return
	}
	if req.Name == "" {
		renderBadRequestErr(w, errors.New("name is null"))
		return
	}
	hospital, err := api.hospitalService.CreateHospital(r.Context(), &req)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	renderJSON(w, http.StatusCreated, hospital)
}

func (api *API) handleGetHospital(w http.ResponseWriter, r *http.Request) {
	hidStr := mux.Vars(r)["id"]
	hid, err := strconv.ParseInt(hidStr, 10, 64)
	if err != nil {
		renderBadRequestErr(w, err)
		return
	}
	hospital, err := api.hospitalService.GetHospital(r.Context(), hid)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	renderJSON(w, http.StatusOK, hospital)
}

func (api *API) handleUpdateHospital(w http.ResponseWriter, r *http.Request) {
	hidStr := mux.Vars(r)["id"]
	hid, err := strconv.ParseInt(hidStr, 10, 64)
	if err != nil {
		renderBadRequestErr(w, err)
		return
	}
	var req dto.Hospital
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		renderBadRequestErr(w, err)
		return
	}
	if req.Name == "" {
		renderBadRequestErr(w, errors.New("name is null"))
		return
	}
	hospital, err := api.hospitalService.GetHospital(r.Context(), hid)
	if err != nil {
		renderSvcError(w, err)
		return
	}
	hospital.Name = req.Name
	hospital.DisplayName = req.DisplayName
	if err := api.hospitalService.UpdateHospital(r.Context(), hospital); err != nil {
		renderSvcError(w, err)
		return
	}
}
