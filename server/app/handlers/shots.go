package handlers

import (
	"encoding/json"
	stderr "errors"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/utils"
	"github.com/cornbuddy/reflectiveTarget/server/domain/errors"
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
	log := utils.LoggerFromCtx(req.Context())

	targetID, err := strconv.Atoi(req.PathValue("targetID"))
	if err != nil {
		log.Error("failed to parse target id", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	log = log.With(zap.Int("target-id", targetID))

	shots, err := h.ShotsDao.List(targetID)
	if stderr.Is(err, errors.ErrNotFound) {
		log.Warn("target not found")
		http.Error(resp, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		log.Error("cannot fetch shots", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(ShotsResponse{Shots: shots})
	if err != nil {
		log.Error("cannot marshal response", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write(data)
}

func (h shotsHandler) post(resp http.ResponseWriter, req *http.Request) {
	log := utils.LoggerFromCtx(req.Context())

	targetID, err := strconv.Atoi(req.PathValue("targetID"))
	if err != nil {
		log.Error("failed to parse target id", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	log = log.With(zap.Int("target-id", targetID))

	var shots ShotsRequest
	if err := json.NewDecoder(req.Body).Decode(&shots); err != nil {
		log.Error("failed to decode body", zap.Error(err))
		http.Error(resp, "bad request", http.StatusBadRequest)
		return
	}

	if res := h.Validator.Validate(shots.Shots); res.IsInvalid() {
		log.Error("failed to decode body",
			zap.Errors("errors", res.Errors),
		)
		http.Error(resp, "invalid payload", http.StatusBadRequest)
		return
	}

	cookies := req.CookiesNamed(constants.SessionCookieName)
	if len(cookies) != 1 {
		log.Error("too much of cookies", zap.Any("cookies", cookies))
		http.Error(resp, "bad cookies", http.StatusBadRequest)
		return
	}

	shooter := cookies[0].Value
	if err := h.ShotsDao.Save(shooter, targetID, shots.Shots); err != nil {
		log.Error("cannot save shots", zap.Error(err))
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("shots saved")
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusCreated)
	resp.Write([]byte("ok"))
}
