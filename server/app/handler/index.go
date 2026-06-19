package handler

import (
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/render"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/utils"
	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
)

type indexHandler struct{}

func (h indexHandler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.URL.Path != "/" {
		log := log.Logger(ctx)
		log.Error("bad url")
		utils.HttpError(w, http.StatusNotFound)

		return
	}

	render.Layout.Index(ctx, w)
}
