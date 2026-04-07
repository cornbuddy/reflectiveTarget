package handler

import (
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/app/render"
)

type indexHandler struct{}

func (h indexHandler) get(resp http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(resp, req)
		return
	}

	render.Layout.Index(req.Context(), resp)
}
