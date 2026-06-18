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
	log := log.Logger(ctx)
	log.Debug("going to parse form")
	if err := r.ParseForm(); err != nil {
		log.Error("failed to parse form")
		utils.HttpError(w, http.StatusBadRequest)
		return
	}

	rawID := mux.Vars(r)["targetID"]
	log.Debug("form parsed, going to parse id from url", zap.String("id", rawID))
	targetID, err := strconv.Atoi(rawID)
	if err != nil {
		log.Warn("failed to parse target id", zap.String("id", rawID), zap.Error(err))
		utils.HttpError(w, http.StatusBadRequest)
		return
	}

	log = log.With(zap.Int("target-id", targetID))
	log.Debug("id parsed, going to parse form", zap.Any("form", r.Form))
	form, err := contracts.NewTargetForm(r.Form)
	if err != nil {
		log.Error("bad form", zap.Error(err))
		utils.HttpError(w, http.StatusBadRequest)
		return
	}

	session := sessiondata.Read(ctx)
	id := valueobjects.ID(targetID)
	log = log.With(zap.Stringer("form", form))
	log.Debug("form parsed, going to build target")
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
	log.Debug("target built, going to save it")
	if err := h.repo.Save(ctx, target); err != nil {
		log.Error("failed to save target", zap.Error(err))
		utils.HttpError(w, http.StatusInternalServerError)
		return
	}

	log.Debug("target is saved")
	utils.Redirect(w, r, "/targets", "target updated")
}

func (h targetsHandler) getExisting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := log.Logger(ctx)
	vars := mux.Vars(r)
	targetID, err := strconv.Atoi(vars["targetID"])
	if err != nil {
		log.Warn("failed to parse target id", zap.Error(err))
		utils.HttpError(w, http.StatusBadRequest)
		return
	}

	log = log.With(zap.Int("target-id", targetID))
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
	log := log.Logger(ctx)
	session := sessiondata.Read(ctx)
	targets, err := h.repo.ListTargetsOfUser(ctx, session.UserID)
	if err != nil {
		log.Error("failed to fetch targets", zap.Error(err))
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
	log := log.Logger(ctx)
	log.Debug("going to parse form")
	if err := r.ParseForm(); err != nil {
		log.Error("failed to parse form", zap.Error(err))
		utils.HttpError(w, http.StatusBadRequest)
		return
	}

	session := sessiondata.Read(ctx)
	log.Debug("form parsed, going to create form DTO", zap.Any("form", r.Form))
	form, err := contracts.NewTargetForm(r.Form)
	if err != nil {
		log.Error("bad form", zap.Any("form", r.Form), zap.Error(err))
		utils.HttpError(w, http.StatusBadRequest)
		return
	}

	log = log.With(zap.Stringer("form", form))
	log.Debug("DTO created, building target object")
	target, err := h.builder.Target(ctx, form, session.UserID, 0)
	if target == nil {
		log.Debug("failed to create target object, form is invalid")
		w.WriteHeader(http.StatusBadRequest)
		view.TargetForm(ctx, w, render.TargetFormData{TargetForm: *form})
		return
	} else if err != nil {
		log.Error("failed to build target object", zap.Error(err))
		utils.HttpError(w, http.StatusInternalServerError)
		return
	}

	log = log.With(zap.Stringer("target", target))
	log.Debug("target is created, saving")
	if err := h.repo.Save(ctx, target); err != nil {
		log.Error("failed to save target", zap.Error(err))
		utils.HttpError(w, http.StatusInternalServerError)
		return
	}

	log.Info("target saved")
	utils.Redirect(w, r, "/targets", "target created")
}
