package handlers

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/forms"
	"github.com/cornbuddy/reflectiveTarget/server/app/render"
	"github.com/cornbuddy/reflectiveTarget/server/app/utils"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

type authzHandler struct {
	daos.UserDao
	daos.SessionStore
	SignupFormValidator
	LoginFormValidator
}

func (h authzHandler) getLogout(w http.ResponseWriter, r *http.Request) {
	log := utils.LoggerFromCtx(r.Context())

	if _, err := utils.SaveSession(h.SessionStore, false, w); err != nil {
		log.Error("failed to save session", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("logout succeeded")
	utils.Redirect(w, r, "/", "Logout succeeded")
}

func (h authzHandler) postLogin(w http.ResponseWriter, r *http.Request) {
	log := utils.LoggerFromCtx(r.Context())

	if err := r.ParseForm(); err != nil {
		log.Error("failed to parse form", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	form := forms.NewLoginForm(r.Form)
	if valid := h.LoginFormValidator.Validate(&form); !valid {
		log.Debug("login failed", zap.Any("form", form))
		w.WriteHeader(http.StatusUnauthorized)
		render.View.Login(r.Context(), w, form)
		return
	}

	if _, err := utils.SaveSession(h.SessionStore, true, w); err != nil {
		log.Error("failed to register session", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("user is logged in",
		zap.String("username", form.Username.Value),
	)
	utils.Redirect(w, r, "/", "Login succeeded")
}

func (h authzHandler) postSignup(w http.ResponseWriter, r *http.Request) {
	log := utils.LoggerFromCtx(r.Context())

	if err := r.ParseForm(); err != nil {
		log.Error("failed to parse form", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	form := forms.NewSignupForm(r.Form)
	if valid := h.SignupFormValidator.Validate(&form); !valid {
		log.Debug("signup failed", zap.Any("form", form))
		w.WriteHeader(http.StatusBadRequest)
		render.View.Signup(r.Context(), w, form)
		return
	}

	username := form.Username.Value
	newUser, err := entities.NewUser(username, form.Password.Value)
	if err != nil {
		log.Error("failed to create user object", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.UserDao.Save(newUser); err != nil {
		log.Error("failed to save user object", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := utils.SaveSession(h.SessionStore, true, w); err != nil {
		log.Error("failed to save session", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("user object created", zap.String("username", username))
	utils.Redirect(w, r, "/", "User successfully created")
}

func (h authzHandler) getSignup(w http.ResponseWriter, r *http.Request) {
	render.Layout.Signup(r.Context(), w, nil)
}

func (h authzHandler) getLogin(w http.ResponseWriter, r *http.Request) {
	render.Layout.Login(r.Context(), w, nil)
}
