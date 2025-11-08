package handlers

import (
	"net/http"
)

func (r IndexRouter) Get(resp http.ResponseWriter, req *http.Request) {
	engine.Render(resp, "views/index.tmpl", nil)
}
