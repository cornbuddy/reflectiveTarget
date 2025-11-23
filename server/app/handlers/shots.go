package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	myerrors "github.com/cornbuddy/reflectiveTarget/server/domain/errors"
	"github.com/cornbuddy/reflectiveTarget/server/domain/validators"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

type ShotsRequest struct {
	valueobjects.Shots `json:"shots"`
}

type ShotsResponse struct {
	valueobjects.Shots `json:"shots"`
}

type shotsHandler struct {
	daos.ShotsDao
	Validator validators.ShotsRequestValidator
}

func (h shotsHandler) get(resp http.ResponseWriter, req *http.Request) {
	targetID, err := strconv.Atoi(req.PathValue("targetID"))
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	shots, err := h.ShotsDao.List(targetID)
	if errors.Is(err, myerrors.ErrNotFound) {
		http.Error(resp, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(ShotsResponse{Shots: shots})
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(data)
}

func (h shotsHandler) post(resp http.ResponseWriter, req *http.Request) {
	var shots ShotsRequest
	if err := json.NewDecoder(req.Body).Decode(&shots); err != nil {
		http.Error(resp, "bad request", http.StatusBadRequest)
		return
	}

	if res := h.Validator.Validate(shots.Shots); res.IsInvalid() {
		http.Error(resp, "invalid payload", http.StatusBadRequest)
		return
	}

	targetID, err := strconv.Atoi(req.PathValue("targetID"))
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	cookies := req.CookiesNamed(constants.SessionCookieName)
	if len(cookies) != 1 {
		http.Error(resp, "bad cookies", http.StatusBadRequest)
		return
	}

	shooter := cookies[0].Value
	if err := h.ShotsDao.Save(shooter, targetID, shots.Shots); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusCreated)
	resp.Write([]byte("ok"))
}
