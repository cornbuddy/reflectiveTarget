package handlers

import (
	"net/http"
)

type indexHandler struct{}

func (h indexHandler) get(resp http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(resp, req)
		return
	}

	engine.Render(resp, "views/index.tmpl", nil)
}
