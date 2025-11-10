package handlers

import (
	"net/http"
)

func (r IndexRouter) Get(resp http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(resp, req)
		return
	}

	engine.Render(resp, "views/index.tmpl", nil)
}
