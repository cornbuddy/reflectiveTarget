package handlers

import (
	"embed"
	"net/http"

	"github.com/abiosoft/mold"
)

type AuthzRouter struct{}

//go:embed templates
var dir embed.FS

var options = mold.With(
	mold.WithRoot("templates"),
	mold.WithLayout("layout.tmpl"),
)
var engine = mold.Must(mold.New(dir, options))

func (r AuthzRouter) GetSignup(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/html")
	engine.Render(resp, "views/signup.tmpl", nil)
}

func (r AuthzRouter) GetLogin(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("Content-Type", "text/html")
	engine.Render(resp, "views/login.tmpl", nil)
}
