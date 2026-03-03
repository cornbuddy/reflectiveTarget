package handlers

import (
	"encoding/json"
	stderr "errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/utils"
	"github.com/cornbuddy/reflectiveTarget/server/domain/errors"
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
	Validator ShotsRequestValidator
}

func (h shotsHandler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := utils.LoggerFromCtx(ctx)

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["targetID"])
	if err != nil {
		log.Error("failed to parse target id", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log = log.With(zap.Int("target-id", targetID))

	shots, err := h.ShotsDao.List(ctx, valueobjects.ID(targetID))
	if stderr.Is(err, errors.ErrNotFound) {
		log.Warn("target not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.Error("cannot fetch shots", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(ShotsResponse{Shots: shots})
	if err != nil {
		log.Error("cannot marshal response", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h shotsHandler) post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := utils.LoggerFromCtx(ctx)

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["targetID"])
	if err != nil {
		log.Error("failed to parse target id", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log = log.With(zap.Int("target-id", targetID))

	var shots ShotsRequest
	if err := json.NewDecoder(r.Body).Decode(&shots); err != nil {
		log.Error("failed to decode body", zap.Error(err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if res := h.Validator.Validate(shots); res.IsInvalid() {
		log.Error("failed to decode body",
			zap.Errors("errors", res.Errors),
		)
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	cookies := r.CookiesNamed(constants.SessionCookieName)
	if len(cookies) != 1 {
		log.Error("too much of cookies", zap.Any("cookies", cookies))
		http.Error(w, "bad cookies", http.StatusBadRequest)
		return
	}

	shooter := cookies[0].Value
	id := valueobjects.ID(targetID)
	if err := h.ShotsDao.Save(ctx, shooter, id, shots.Shots); err != nil {
		log.Error("cannot save shots", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("shots saved")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("ok"))
}
