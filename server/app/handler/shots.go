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
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/utils"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/validators"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
)

type shotsHandler struct {
	dao       daos.ShotsDao
	validator validators.ShotsRequestValidator
}

func (h shotsHandler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := log.Logger(ctx)

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["targetID"])
	if err != nil {
		log.Warn("failed to parse target id")
		http.Error(w, "bad target id", http.StatusBadRequest)

		return
	}

	log = log.With(zap.Int("target-id", targetID))
	shots, err := h.dao.List(ctx, valueobjects.ID(targetID))
	if errors.Is(err, sql.ErrNoRows) {
		log.Warn("target not found")
		http.Error(w, "target not found", http.StatusNotFound)

		return
	} else if err != nil {
		log.Error("could not fetch shots", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	data, err := json.Marshal(contracts.ShotsResponse{Shots: shots})
	if err != nil {
		log.Error("could unmarshal shots", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	utils.LogBadWrites(log)(w.Write(data))
}

func (h shotsHandler) post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := log.Logger(ctx)

	vars := mux.Vars(r)
	rawID := vars["targetID"]
	targetID, err := strconv.Atoi(rawID)
	if err != nil {
		log.Warn("failed to parse target id", zap.String("target-id", rawID))
		http.Error(w, "bad target id", http.StatusBadRequest)

		return
	}

	log = log.With(zap.Int("target-id", targetID))
	var shots contracts.ShotsRequest
	if err := json.NewDecoder(r.Body).Decode(&shots); err != nil {
		log.Warn("failed to decode body")
		http.Error(w, "bad request body", http.StatusBadRequest)

		return
	}

	if res := h.validator.Validate(shots); res.IsInvalid() {
		log.Warn("invalid shots", zap.Any("shots", shots))
		http.Error(w, "bad request body", http.StatusBadRequest)

		return
	}

	cookies := r.CookiesNamed(constants.SessionCookieName)
	if len(cookies) != 1 {
		log.Warn("too much cookies", zap.Any("cookies", cookies))
		http.Error(w, "bad request", http.StatusBadRequest)

		return
	}

	shooter := cookies[0].Value
	id := valueobjects.ID(targetID)
	if err := h.dao.Save(ctx, shooter, id, shots.Shots); err != nil {
		log.Error("failed to save shots", zap.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	log.Info("shots saved")
	w.WriteHeader(http.StatusCreated)
}
