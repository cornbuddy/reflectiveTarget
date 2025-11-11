package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/private/daos"
	"github.com/cornbuddy/reflectiveTarget/server/private/model"
	"github.com/cornbuddy/reflectiveTarget/server/private/validators"
)

type ShotsRouter struct {
	daos.ShotsDao
	Validator validators.ShotsRequestValidator
}

func (r ShotsRouter) Post(resp http.ResponseWriter, req *http.Request) {
	var shots model.ShotsRequest
	if err := json.NewDecoder(req.Body).Decode(&shots); err != nil {
		http.Error(resp, "bad request", http.StatusBadRequest)
		return
	}

	if res := r.Validator.Validate(shots); res.IsInvalid() {
		http.Error(resp, "invalid payload", http.StatusBadRequest)
		return
	}

	cookies := req.CookiesNamed(SessionCookieName)
	if len(cookies) != 1 {
		http.Error(resp, "bad cookies", http.StatusBadRequest)
	}

	if err := r.ShotsDao.Save(shots.Shots); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusCreated)
}
