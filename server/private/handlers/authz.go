package handlers

import (
	"embed"
	"net/http"

	"github.com/abiosoft/mold"
	"github.com/cornbuddy/reflectiveTarget/server/private/daos"
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
	dao := r.UserDao
	req.ParseForm()
	user, err := dao.Save(user)
}

func (r AuthzRouter) GetSignup(resp http.ResponseWriter, req *http.Request) {
	engine.Render(resp, "views/signup.tmpl", nil)
}

func (r AuthzRouter) GetLogin(resp http.ResponseWriter, req *http.Request) {
	engine.Render(resp, "views/login.tmpl", nil)
}
