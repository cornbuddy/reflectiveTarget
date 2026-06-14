package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/builders"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/render"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/utils"
	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
	"github.com/cornbuddy/reflectiveTarget/server/infra/repositories"
)

type targetsHandler struct {
	repo    repositories.TargetRepo
	builder builders.TargetBuilder
}

func (h targetsHandler) putExisting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		log.Error(ctx, "failed to parse form")
		utils.HttpError(w, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["targetID"])
	if err != nil {
		log.Warn(ctx, "failed to parse target id",
			zap.String("id", vars["targetID"]), zap.Error(err))
		utils.HttpError(w, http.StatusBadRequest)
		return
	}

	log := log.Logger(ctx).With(zap.Int("target-id", targetID))
	form, err := contracts.NewTargetForm(r.Form)
	if err != nil {
		log.Error("bad form", zap.Error(err))
		utils.HttpError(w, http.StatusBadRequest)
		return
	}

	log = log.With(zap.Stringer("form", form))
	session := sessiondata.Read(ctx)
	id := valueobjects.ID(targetID)
	target, err := h.builder.Target(ctx, form, session.UserID, id)
	if err != nil {
		log.Error("failed to build target", zap.Error(err))
		utils.HttpError(w, http.StatusInternalServerError)
		return
	} else if target == nil {
		log.Debug("form is invalid")
		w.WriteHeader(http.StatusBadRequest)
		view.TargetForm(ctx, w, render.TargetFormData{
			TargetForm: *form,
			ID:         id,
		})
		return
	}

	log = log.With(zap.Stringer("target", target))
	target.ID = id
	if err := h.repo.Save(ctx, target); err != nil {
		log.Error("failed to save target", zap.Error(err))
		utils.HttpError(w, http.StatusInternalServerError)
		return
	}

	utils.Redirect(w, r, "/targets", "target updated")
}

func (h targetsHandler) getExisting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["targetID"])
	if err != nil {
		log.Warn(ctx, "failed to parse target id", zap.Error(err))
		utils.HttpError(w, http.StatusBadRequest)
		return
	}

	log := log.Logger(ctx).With(zap.Int("target-id", targetID))
	target, err := h.repo.Get(ctx, valueobjects.ID(targetID))
	if err != nil {
		log.Error("failed to fetch target", zap.Error(err))
		utils.HttpError(w, http.StatusInternalServerError)
		return
	} else if target == nil {
		log.Error("target not found", zap.Error(err))
		utils.HttpError(w, http.StatusNotFound)
		return
	}

	layout.TargetForm(ctx, w, render.TargetFormData{
		TargetForm: contracts.NewTargetFormFromTarget(*target),
		ID:         target.ID,
	})
}

func (h targetsHandler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := sessiondata.Read(ctx)
	targets, err := h.repo.ListTargetsOfUser(ctx, session.UserID)
	if err != nil {
		log.Error(ctx, "failed to fetch targets", zap.Error(err))
		utils.HttpError(w, http.StatusInternalServerError)
		return
	}

	layout.Targets(ctx, w, render.TargetsData{Targets: targets})
}

func (h targetsHandler) getNew(w http.ResponseWriter, r *http.Request) {
	layout.TargetForm(r.Context(), w, render.TargetFormData{})
}

func (h targetsHandler) postNew(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		log.Error(ctx, "failed to parse form", zap.Error(err))
		utils.HttpError(w, http.StatusBadRequest)
		return
	}

	session := sessiondata.Read(ctx)
	form, err := contracts.NewTargetForm(r.Form)
	if err != nil {
		log.Error(ctx, "bad form", zap.Error(err))
		utils.HttpError(w, http.StatusBadRequest)
		return
	}

	target, err := h.builder.Target(ctx, form, session.UserID, 0)
	if target == nil {
		log.Debug(ctx, "form is invalid", zap.Any("form", form))
		w.WriteHeader(http.StatusBadRequest)
		view.TargetForm(ctx, w, render.TargetFormData{TargetForm: *form})
		return
	} else if err != nil {
		log.Error(ctx, "failed to validate form", zap.Error(err))
		utils.HttpError(w, http.StatusInternalServerError)
		return
	}

	if err := h.repo.Save(ctx, target); err != nil {
		log.Error(ctx, "failed to save target", zap.Error(err))
		utils.HttpError(w, http.StatusInternalServerError)
		return
	}

	utils.Redirect(w, r, "/targets", "target created")
}
