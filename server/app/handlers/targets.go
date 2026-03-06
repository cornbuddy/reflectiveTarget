package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/render"
	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/infra/repositories"
)

type targetsHandler struct {
	repo repositories.TargetRepo
}

func (h targetsHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := sessiondata.Read(ctx)
	log := session.Logger

	targets, err := h.repo.ListTargetNamesOfUser(ctx, session.Username)
	if err != nil {
		msg := "failed to fetch list of targets"
		log.Error(msg, zap.Error(err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	render.Layout.Targets(ctx, w, targets)
}

func (h targetsHandler) new(w http.ResponseWriter, r *http.Request) {}

func (h targetsHandler) update(w http.ResponseWriter, r *http.Request) {}
