package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/validators"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
)

type shotsHandler struct {
	daos.ShotsDao
	Validator validators.ShotsRequestValidator
}

func (h shotsHandler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.Logger(ctx)

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["targetID"])
	if err != nil {
		logger.Warn("failed to parse target id")
		http.Error(w, "bad target id", http.StatusBadRequest)
		return
	}

	logger = logger.With(zap.Int("target-id", targetID))
	shots, err := h.ShotsDao.List(ctx, valueobjects.ID(targetID))
	if errors.Is(err, sql.ErrNoRows) {
		logger.Warn("target not found")
		http.Error(w, "target not found", http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error("could not fetch shots", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(contracts.ShotsResponse{Shots: shots})
	if err != nil {
		logger.Error("could unmarshal shots", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h shotsHandler) post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := log.Logger(ctx)

	vars := mux.Vars(r)
	rawID := vars["targetID"]
	targetID, err := strconv.Atoi(rawID)
	if err != nil {
		logger.Warn("failed to parse target id", zap.String("target-id", rawID))
		http.Error(w, "bad target id", http.StatusBadRequest)
		return
	}

	logger = logger.With(zap.Int("target-id", targetID))
	var shots contracts.ShotsRequest
	if err := json.NewDecoder(r.Body).Decode(&shots); err != nil {
		logger.Warn("failed to decode body")
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	if res := h.Validator.Validate(shots); res.IsInvalid() {
		logger.Warn("invalid shots", zap.Any("shots", shots))
		http.Error(w, "bad request body", http.StatusBadRequest)
		return
	}

	cookies := r.CookiesNamed(constants.SessionCookieName)
	if len(cookies) != 1 {
		logger.Warn("too much cookies", zap.Any("cookies", cookies))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	shooter := cookies[0].Value
	id := valueobjects.ID(targetID)
	if err := h.ShotsDao.Save(ctx, shooter, id, shots.Shots); err != nil {
		logger.Error("failed to save shots", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info("shots saved")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("ok"))
}
