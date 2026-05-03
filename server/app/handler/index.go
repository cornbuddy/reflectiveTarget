package handler

import (
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/render"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/utils"
)

type indexHandler struct{}

func (h indexHandler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.URL.Path != "/" {
		log := utils.LoggerFromCtx(ctx)
		utils.NotFound(log, w, "not found", nil)
		return
	}

	render.Layout.Index(ctx, w)
}
