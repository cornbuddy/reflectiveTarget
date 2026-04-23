package handler

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/render"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/validators"
	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/app/utils"
	"github.com/cornbuddy/reflectiveTarget/server/infra/repositories"
)

type targetsHandler struct {
	repo          repositories.TargetRepo
	formValidator validators.TargetFormValidator
}

func (h targetsHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := sessiondata.Read(ctx)
	log := utils.LoggerFromCtx(ctx)

	targets, err := h.repo.ListTargetsOfUser(ctx, session.Username)
	if err != nil {
		msg := "failed to fetch list of targets"
		log.Error(msg, zap.Error(err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	render.Layout.Targets(ctx, w, render.TargetsData{Targets: targets})
}

func (h targetsHandler) makeNew(w http.ResponseWriter, r *http.Request) {
	render.Layout.TargetForm(r.Context(), w, render.TargetFormData{})
}

func (h targetsHandler) saveNew(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := utils.LoggerFromCtx(ctx)

	if err := r.ParseForm(); err != nil {
		log.Error("failed to parse form", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	form := contracts.NewTargetForm(r.Form)
	if valid := h.formValidator.Validate(&form); !valid {
		log.Debug("form is invalid", zap.Any("form", form))
		w.WriteHeader(http.StatusBadRequest)
		render.View.TargetForm(ctx, w, render.TargetFormData{TargetForm: form})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h targetsHandler) update(w http.ResponseWriter, r *http.Request) {}
