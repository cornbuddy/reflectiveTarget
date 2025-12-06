package render

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/utils"
)

func (e engine) Signup(ctx context.Context, w http.ResponseWriter, data any) {
	e.render(ctx, "views/signup.tmpl", w, data)
}

func (e engine) Login(ctx context.Context, w http.ResponseWriter, data any) {
	e.render(ctx, "views/login.tmpl", w, data)
}

func (e engine) Index(ctx context.Context, w http.ResponseWriter, data any) {
	e.render(ctx, "views/index.tmpl", w, data)
}

func (e engine) render(
	ctx context.Context, path string, w http.ResponseWriter, data any,
) {
	ctxData := extractDataFromContext(ctx)
	resultData := viewData{ctxData, data}
	if err := e.Engine.Render(w, path, resultData); err != nil {
		log := utils.LoggerFromCtx(ctx)
		log.Error("failed to render",
			zap.Error(err),
			zap.String("path", path),
		)
	}

}

func extractDataFromContext(ctx context.Context) contextData {
	isAuthorized := false
	auth := ctx.Value(constants.AuthenticatedCtx)
	if auth != nil {
		isAuthorized = *auth.(*bool)
	}

	return contextData{
		IsAuthorized: isAuthorized,
	}
}
