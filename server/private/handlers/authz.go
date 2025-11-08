package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/cornbuddy/reflectiveTarget/server/private/model"
)

const SessionCookieName = "session-token"

func (r AuthzRouter) PostLogin(resp http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	userObj, err := model.NewUser(req.Form)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := r.UserDao.Find(userObj.Username)
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
		resp.Write([]byte("Login succeeded"))
	} else {
		http.Error(resp, "Password is wrong", http.StatusUnauthorized)
	}

}

func (r AuthzRouter) PostSignup(resp http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	newUser, err := model.NewUser(req.Form)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := r.UserDao.Find(newUser.Username)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	} else if user != nil {
		http.Error(resp, "User already exists", http.StatusConflict)
		return
	}

	if err := r.UserDao.Save(*newUser); err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(resp, &http.Cookie{
		Name:    SessionCookieName,
		Value:   uuid.NewString(),
		Expires: time.Now().AddDate(0, 1, 0),
	})
	resp.WriteHeader(http.StatusCreated)
	resp.Write([]byte("User successfully created"))
}

func (r AuthzRouter) GetSignup(resp http.ResponseWriter, req *http.Request) {
	engine.Render(resp, "views/signup.tmpl", nil)
}

func (r AuthzRouter) GetLogin(resp http.ResponseWriter, req *http.Request) {
	engine.Render(resp, "views/login.tmpl", nil)
}
