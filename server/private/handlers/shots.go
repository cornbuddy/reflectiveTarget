package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cornbuddy/reflectiveTarget/server/private/daos"
	"github.com/cornbuddy/reflectiveTarget/server/private/model"
	"github.com/cornbuddy/reflectiveTarget/server/private/validators"
)

type shotsHandler struct {
	daos.ShotsDao
	Validator validators.ShotsRequestValidator
}

func (h shotsHandler) post(resp http.ResponseWriter, req *http.Request) {
	var shots model.ShotsRequest
	if err := json.NewDecoder(req.Body).Decode(&shots); err != nil {
		http.Error(resp, "bad request", http.StatusBadRequest)
		return
	}

	if res := h.Validator.Validate(shots); res.IsInvalid() {
		http.Error(resp, "invalid payload", http.StatusBadRequest)
		return
	}

	targetID, err := strconv.Atoi(req.PathValue("targetID"))
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	cookies := req.CookiesNamed(SessionCookieName)
	if len(cookies) != 1 {
		http.Error(resp, "bad cookies", http.StatusBadRequest)
		return
	}

	shooter := cookies[0].Value
	if err := h.ShotsDao.Save(shooter, targetID, shots.Shots); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusCreated)
}
