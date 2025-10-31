package handlers

import (
	"embed"
	"net/http"

	"github.com/abiosoft/mold"

	"github.com/cornbuddy/reflectiveTarget/server/private/daos"
	"github.com/cornbuddy/reflectiveTarget/server/private/model"
)

type AuthzRouter struct {
	daos.UserDao
}

//go:embed templates
var dir embed.FS
var options = mold.With(
	mold.WithRoot("templates"),
	mold.WithLayout("layout.tmpl"),
)
var engine = mold.Must(mold.New(dir, options))

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

	resp.WriteHeader(http.StatusCreated)
	resp.Write([]byte("User successfully created"))
}

func (r AuthzRouter) GetSignup(resp http.ResponseWriter, req *http.Request) {
	engine.Render(resp, "views/signup.tmpl", nil)
}

func (r AuthzRouter) GetLogin(resp http.ResponseWriter, req *http.Request) {
	engine.Render(resp, "views/login.tmpl", nil)
}
