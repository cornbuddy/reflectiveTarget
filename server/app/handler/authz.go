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
	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
)

type authzHandler struct {
	userDao         daos.UserDao
	sessionStore    daos.SessionStore
	signupValidator validators.SignupFormValidator
	loginValidator  validators.LoginFormValidator
}

func (h authzHandler) getLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := log.Logger(ctx)
	data := sessiondata.SessionData{}
	if _, err := utils.SaveSession(ctx, h.sessionStore, data, w); err != nil {
		log.Error("failed to save session")
		utils.HttpError(w, http.StatusInternalServerError)

		return
	}

	log.Info("logout succeeded")
	utils.Redirect(w, r, "/", "Logout succeeded")
}

func (h authzHandler) postLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := log.Logger(ctx)
	if err := r.ParseForm(); err != nil {
		log.Warn("failed to parse form", zap.Error(err))
		utils.HttpError(w, http.StatusBadRequest)

		return
	}

	form := contracts.NewLoginForm(r.Form)
	if valid := h.loginValidator.Validate(ctx, &form); !valid {
		log.Debug("login failed", zap.Any("form", form))
		w.WriteHeader(http.StatusUnauthorized)
		render.View.Login(ctx, w, render.LoginData{LoginForm: form})

		return
	}

	user, err := h.userDao.Find(ctx, form.Username.Value)
	if err != nil {
		log.Error("failed to fetch user", zap.Error(err))
		utils.HttpError(w, http.StatusInternalServerError)

		return
	}

	log = log.With(zap.String("username", user.Username))
	data := sessiondata.SessionData{
		IsAuthenticated: true,
		UserID:          user.ID,
		Username:        user.Username,
	}
	if _, err := utils.SaveSession(ctx, h.sessionStore, data, w); err != nil {
		log.Error("failed to register session", zap.Error(err))
		utils.HttpError(w, http.StatusInternalServerError)

		return
	}

	log.Info("user is logged in")
	utils.Redirect(w, r, "/", "Login succeeded")
}

func (h authzHandler) postSignup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := log.Logger(ctx)
	if err := r.ParseForm(); err != nil {
		log.Warn("failed to parse form", zap.Error(err))
		utils.HttpError(w, http.StatusBadRequest)

		return
	}

	form := contracts.NewSignupForm(r.Form)
	if valid, err := h.signupValidator.Validate(ctx, &form); err != nil {
		log.Error("failed to validate form", zap.Error(err))
		utils.HttpError(w, http.StatusInternalServerError)

		return
	} else if !valid {
		log.Debug("signup failed", zap.Any("form", form))
		w.WriteHeader(http.StatusBadRequest)
		render.View.Signup(ctx, w, render.SignupData{SignupForm: form})

		return
	}

	username := form.Username.Value
	log = log.With(zap.String("username", username))
	user, err := entities.NewUser(username, form.Password.Value)
	if err != nil {
		log.Error("failed to create user object", zap.Error(err))
		utils.HttpError(w, http.StatusInternalServerError)

		return
	}

	if err := h.userDao.Save(ctx, user); err != nil {
		log.Error("failed to save user object", zap.Error(err))
		utils.HttpError(w, http.StatusInternalServerError)

		return
	}

	data := sessiondata.SessionData{
		IsAuthenticated: true,
		UserID:          user.ID,
		Username:        user.Username,
	}
	store := h.sessionStore
	if _, err := utils.SaveSession(ctx, store, data, w); err != nil {
		log.Error("failed to save session", zap.Error(err))
		utils.HttpError(w, http.StatusInternalServerError)

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
