package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/cornbuddy/reflectiveTarget/server/domain/constants"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

const SessionCookieName = "session-token"

type authzHandler struct {
	daos.UserDao
}

func (h authzHandler) postLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   SessionCookieName,
		Value:  "",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
	w.Write([]byte("Logout succeeded"))
}

func (h authzHandler) postLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userObj, err := entities.NewUser(r.Form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.UserDao.Find(userObj.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if user == nil {
		http.Error(w, "User does not exist", http.StatusUnauthorized)
		return
	}

	authorized, err := user.Password.Verify(r.Form.Get("password"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if authorized {
		http.SetCookie(w, &http.Cookie{
			Name:    SessionCookieName,
			Value:   uuid.NewString(),
			Expires: time.Now().Add(constants.SessionDuration),
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
		w.Write([]byte("Login succeeded"))
	} else {
		http.Error(w, "Password is wrong", http.StatusUnauthorized)
	}

}

func (h authzHandler) postSignup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser, err := entities.NewUser(r.Form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.UserDao.Find(newUser.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if user != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	if err := h.UserDao.Save(newUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    SessionCookieName,
		Value:   uuid.NewString(),
		Expires: time.Now().Add(constants.SessionDuration),
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
	w.Write([]byte("User successfully created"))
}

func (h authzHandler) getSignup(w http.ResponseWriter, r *http.Request) {
	layout.Render(w, "views/signup.tmpl", nil)
}

func (h authzHandler) getLogin(w http.ResponseWriter, r *http.Request) {
	layout.Render(w, "views/login.tmpl", nil)
}
