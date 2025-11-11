package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cornbuddy/reflectiveTarget/server/private/model"
)

func (r ShotsHandler) Post(resp http.ResponseWriter, req *http.Request) {
	var shots model.ShotsRequest
	if err := json.NewDecoder(req.Body).Decode(&shots); err != nil {
		http.Error(resp, "bad request", http.StatusBadRequest)
		return
	}

	if res := r.Validator.Validate(shots); res.IsInvalid() {
		http.Error(resp, "invalid payload", http.StatusBadRequest)
		return
	}

	// path value doesn't work out of the box, because it needs to be
	// registered via mux.HandleFunc at helpers. it's possible to group
	// routes like here
	// https://dev.to/kengowada/go-routing-101-handling-and-grouping-routes-with-nethttp-4k0e
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
	if err := r.ShotsDao.Save(shooter, targetID, shots.Shots); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusCreated)
}
