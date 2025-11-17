package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
	"github.com/cornbuddy/reflectiveTarget/server/model/entities"
)

const SessionCookieName = "session-token"

type authzHandler struct {
	daos.UserDao
}

func (h authzHandler) postLogout(resp http.ResponseWriter, req *http.Request) {
	http.SetCookie(resp, &http.Cookie{
		Name:   SessionCookieName,
		Value:  "",
		MaxAge: -1,
	})
	http.Redirect(resp, req, "/", http.StatusSeeOther)
	resp.Write([]byte("Logout succeeded"))
}

func (h authzHandler) postLogin(resp http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	userObj, err := entities.NewUser(req.Form)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.UserDao.Find(userObj.Username)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	} else if user == nil {
		http.Error(resp, "User does not exist", http.StatusUnauthorized)
		return
	}

	authorized, err := user.Password.Verify(req.Form.Get("password"))
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if authorized {
		http.SetCookie(resp, &http.Cookie{
			Name:    SessionCookieName,
			Value:   uuid.NewString(),
			Expires: time.Now().AddDate(0, 1, 0),
		})
		http.Redirect(resp, req, "/", http.StatusSeeOther)
		resp.Write([]byte("Login succeeded"))
	} else {
		http.Error(resp, "Password is wrong", http.StatusUnauthorized)
	}

}

func (h authzHandler) postSignup(resp http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	newUser, err := entities.NewUser(req.Form)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.UserDao.Find(newUser.Username)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	} else if user != nil {
		http.Error(resp, "User already exists", http.StatusConflict)
		return
	}

	if err := h.UserDao.Save(newUser); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(resp, &http.Cookie{
		Name:    SessionCookieName,
		Value:   uuid.NewString(),
		Expires: time.Now().AddDate(0, 1, 0),
	})
	http.Redirect(resp, req, "/", http.StatusSeeOther)
	resp.Write([]byte("User successfully created"))
}

func (h authzHandler) getSignup(resp http.ResponseWriter, req *http.Request) {
	engine.Render(resp, "views/signup.tmpl", nil)
}

func (h authzHandler) getLogin(resp http.ResponseWriter, req *http.Request) {
	engine.Render(resp, "views/login.tmpl", nil)
}
