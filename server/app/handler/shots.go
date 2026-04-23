package handler

import (
	"encoding/json"
	stderr "errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/utils"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/validators"
	"github.com/cornbuddy/reflectiveTarget/server/domain/errors"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

type shotsHandler struct {
	daos.ShotsDao
	Validator validators.ShotsRequestValidator
}

func (h shotsHandler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := utils.LoggerFromCtx(ctx)

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["targetID"])
	if err != nil {
		utils.BadRequest(log, w, "failed to parse target id", err)
		return
	}

	log = log.With(zap.Int("target-id", targetID))

	shots, err := h.ShotsDao.List(ctx, valueobjects.ID(targetID))
	if stderr.Is(err, errors.ErrNotFound) {
		utils.NotFound(log, w, "target not found", err)
		return
	} else if err != nil {
		utils.InternalServerError(log, w, "cannot fetch shots", err)
		return
	}

	data, err := json.Marshal(contracts.ShotsResponse{Shots: shots})
	if err != nil {
		utils.InternalServerError(log, w, "cannot marshal response", err)
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
		utils.BadRequest(log, w, "failed to parse target id", err)
		return
	}

	log = log.With(zap.Int("target-id", targetID))

	var shots contracts.ShotsRequest
	if err := json.NewDecoder(r.Body).Decode(&shots); err != nil {
		utils.BadRequest(log, w, "failed to decode body", err)
		return
	}

	if res := h.Validator.Validate(shots); res.IsInvalid() {
		utils.BadRequest(log, w, "invalid payload", res.Errors)
		return
	}

	cookies := r.CookiesNamed(constants.SessionCookieName)
	if len(cookies) != 1 {
		utils.BadRequest(
			log.With(zap.Any("cookies", cookies)),
			w, "too much cookies", nil,
		)
		return
	}

	shooter := cookies[0].Value
	id := valueobjects.ID(targetID)
	if err := h.ShotsDao.Save(ctx, shooter, id, shots.Shots); err != nil {
		utils.InternalServerError(log, w, "cannot save shots", err)
		return
	}

	log.Info("shots saved")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("ok"))
}
