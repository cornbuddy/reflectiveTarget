package handlers

import (
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/utils"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

type authzHandler struct {
	daos.UserDao
	daos.SessionStore
}

func (h authzHandler) postLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   constants.SessionCookieName,
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
		_, err := utils.SaveSession(h.SessionStore, true, w)
		if err != nil {
			msg := err.Error()
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

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

	if _, err := utils.SaveSession(h.SessionStore, true, w); err != nil {
		msg := err.Error()
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	w.Write([]byte("User successfully created"))
}

func (h authzHandler) getSignup(w http.ResponseWriter, r *http.Request) {
	layout.Render(w, "views/signup.tmpl", nil)
}

func (h authzHandler) getLogin(w http.ResponseWriter, r *http.Request) {
	layout.Render(w, "views/login.tmpl", nil)
}
