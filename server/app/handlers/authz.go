package handlers

import (
	"net/http"

	"go.uber.org/zap"

	. "github.com/cornbuddy/reflectiveTarget/server/app/logger"
	"github.com/cornbuddy/reflectiveTarget/server/app/utils"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

type authzHandler struct {
	daos.UserDao
	daos.SessionStore
}

func (h authzHandler) postLogout(w http.ResponseWriter, r *http.Request) {
	if _, err := utils.SaveSession(h.SessionStore, false, w); err != nil {
		Log.Error("failed to save session", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Log.Info("logout succeeded")
	http.Redirect(w, r, "/", http.StatusSeeOther)
	w.Write([]byte("Logout succeeded"))
}

func (h authzHandler) postLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		Log.Error("failed to parse form", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userObj, err := entities.NewUser(r.Form)
	if err != nil {
		Log.Error("failed to create user object", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := userObj.Username
	user, err := h.UserDao.Find(username)
	if err != nil {
		Log.Error("could not fetch user object", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if user == nil {
		Log.Warn("requested user not found")
		http.Error(w, "User does not exist", http.StatusUnauthorized)
		return
	}

	authorized, err := user.Password.Verify(r.Form.Get("password"))
	if err != nil {
		Log.Error("failed to verify password", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if authorized {
		_, err := utils.SaveSession(h.SessionStore, true, w)
		if err != nil {
			Log.Error("failed to register session", zap.Error(err))
			msg := err.Error()
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		Log.Info("user is logged in", zap.String("username", username))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		w.Write([]byte("Login succeeded"))
	} else {
		http.Error(w, "Password is wrong", http.StatusUnauthorized)
	}

}

func (h authzHandler) postSignup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		Log.Error("failed to parse form", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser, err := entities.NewUser(r.Form)
	if err != nil {
		Log.Error("failed to create user object", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := newUser.Username
	user, err := h.UserDao.Find(newUser.Username)
	if err != nil {
		Log.Error("could not fetch user object", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if user != nil {
		Log.Warn("username is taken", zap.String("username", username))
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	if err := h.UserDao.Save(newUser); err != nil {
		Log.Error("failed to save user object", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := utils.SaveSession(h.SessionStore, true, w); err != nil {
		Log.Error("failed to save session", zap.Error(err))
		msg := err.Error()
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	Log.Info("user object created", zap.String("username", username))
	http.Redirect(w, r, "/", http.StatusSeeOther)
	w.Write([]byte("User successfully created"))
}

func (h authzHandler) getSignup(w http.ResponseWriter, r *http.Request) {
	layout.Render(w, "views/signup.tmpl", nil)
}

func (h authzHandler) getLogin(w http.ResponseWriter, r *http.Request) {
	layout.Render(w, "views/login.tmpl", nil)
}
