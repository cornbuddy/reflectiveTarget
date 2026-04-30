package handler

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/builders"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/render"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/utils"
	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/infra/repositories"
)

type targetsHandler struct {
	repo    repositories.TargetRepo
	builder builders.TargetBuilder
}

func (h targetsHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := sessiondata.Read(ctx)
	log := utils.LoggerFromCtx(ctx)

	targets, err := h.repo.ListTargetsOfUser(ctx, session.Username)
	if err != nil {
		utils.InternalServerError(log, w, "failed to fetch list of targets", err)
		return
	}

	layout.Targets(ctx, w, render.TargetsData{Targets: targets})
}

func (h targetsHandler) makeNew(w http.ResponseWriter, r *http.Request) {
	layout.TargetForm(r.Context(), w, render.TargetFormData{})
}

func (h targetsHandler) saveNew(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := utils.LoggerFromCtx(ctx)

	if err := r.ParseForm(); err != nil {
		utils.BadRequest(log, w, "failed to parse form", err)
		return
	}

	session := sessiondata.Read(ctx)
	form := contracts.NewTargetForm(r.Form)
	target := h.builder.Target(&form, session.UserID)
	if target == nil {
		log.Debug("form is invalid", zap.Any("form", form))
		w.WriteHeader(http.StatusBadRequest)
		view.TargetForm(ctx, w, render.TargetFormData{TargetForm: form})
		return
	}

	if err := h.repo.Save(ctx, target); err != nil {
		utils.InternalServerError(log, w, "failed to save target", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h targetsHandler) update(w http.ResponseWriter, r *http.Request) {}
