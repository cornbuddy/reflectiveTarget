package handler

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/render"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/utils"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/validators"
	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

type authzHandler struct {
	daos.UserDao
	daos.SessionStore
	validators.SignupFormValidator
	validators.LoginFormValidator
}

func (h authzHandler) getLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := utils.LoggerFromCtx(ctx)

	store := h.SessionStore
	data := sessiondata.SessionData{}
	if _, err := utils.SaveSession(ctx, store, data, w); err != nil {
		utils.InternalServerError(log, w, "failed to save session", err)
		return
	}

	log.Info("logout succeeded")
	utils.Redirect(w, r, "/", "Logout succeeded")
}

func (h authzHandler) postLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := utils.LoggerFromCtx(ctx)

	if err := r.ParseForm(); err != nil {
		utils.BadRequest(log, w, "failed to parse form", err)
		return
	}

	form := contracts.NewLoginForm(r.Form)
	if valid := h.LoginFormValidator.Validate(ctx, &form); !valid {
		log.Debug("login failed", zap.Any("form", form))
		w.WriteHeader(http.StatusUnauthorized)
		render.View.Login(ctx, w, render.LoginData{LoginForm: form})
		return
	}

	user, err := h.UserDao.Find(ctx, form.Username.Value)
	if err != nil {
		utils.InternalServerError(log, w, "failed to fetch user", err)
		return
	}

	log = log.With(zap.String("username", user.Username))
	data := sessiondata.SessionData{
		IsAuthenticated: true,
		UserID:          user.ID,
		Username:        user.Username,
	}
	store := h.SessionStore
	if _, err := utils.SaveSession(ctx, store, data, w); err != nil {
		utils.InternalServerError(log, w, "failed to register session", err)
		return
	}

	log.Info("user is logged in")
	utils.Redirect(w, r, "/", "Login succeeded")
}

func (h authzHandler) postSignup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := utils.LoggerFromCtx(ctx)

	if err := r.ParseForm(); err != nil {
		utils.BadRequest(log, w, "failed to parse form", err)
		return
	}

	form := contracts.NewSignupForm(r.Form)
	if valid, err := h.SignupFormValidator.Validate(ctx, &form); err != nil {
		utils.InternalServerError(log, w, "failed to validate form", err)
		return
	} else if !valid {
		log.Debug("signup failed", zap.Any("form", form))
		w.WriteHeader(http.StatusBadRequest)
		render.View.Signup(ctx, w, render.SignupData{SignupForm: form})
		return
	}

	username := form.Username.Value
	user, err := entities.NewUser(username, form.Password.Value)
	if err != nil {
		utils.InternalServerError(log, w, "failed to create user object", err)
		return
	}

	if err := h.UserDao.Save(ctx, user); err != nil {
		utils.InternalServerError(log, w, "failed to save user object", err)
		return
	}

	data := sessiondata.SessionData{
		IsAuthenticated: true,
		UserID:          user.ID,
		Username:        user.Username,
	}
	store := h.SessionStore
	if _, err := utils.SaveSession(ctx, store, data, w); err != nil {
		utils.InternalServerError(log, w, "failed to save session", err)
		return
	}

	log.Info("user object created", zap.String("username", username))
	utils.Redirect(w, r, "/", "User successfully created")
}

func (h authzHandler) getSignup(w http.ResponseWriter, r *http.Request) {
	render.Layout.Signup(r.Context(), w, render.SignupData{})
}

func (h authzHandler) getLogin(w http.ResponseWriter, r *http.Request) {
	render.Layout.Login(r.Context(), w, render.LoginData{})
}
