package render

import (
	"context"
	"net/http"

	"github.com/go-viper/mapstructure/v2"
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
	log := utils.LoggerFromCtx(ctx).With(
		zap.Any("data", data),
		zap.String("path", path),
	)

	viewData, err := makeViewData(ctx, data)
	if err != nil {
		log.Error("failed to make view data", zap.Error(err))
	}

	if err := e.Engine.Render(w, path, viewData); err != nil {
		log.Error("failed to render", zap.Error(err))
	}

}

func makeViewData(ctx context.Context, data any) (viewData, error) {
	var result viewData
	ctxData := extractDataFromContext(ctx)
	for _, data := range []any{ctxData, data} {
		if err := mapstructure.Decode(data, &result); err != nil {
			return nil, err
		}
	}

	return result, nil
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
