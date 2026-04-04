package render

import (
	"context"
	"net/http"

	"github.com/go-viper/mapstructure/v2"
	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/app/utils"
)

func (e Engine) Targets(
	ctx context.Context, w http.ResponseWriter, data TargetsData,
) {
	e.render(ctx, "views/targets.tmpl", w, data)
}

func (e Engine) Signup(
	ctx context.Context, w http.ResponseWriter, data SignupData,
) {
	e.render(ctx, "views/signup.tmpl", w, data)
}

func (e Engine) Login(
	ctx context.Context, w http.ResponseWriter, data LoginData,
) {
	e.render(ctx, "views/login.tmpl", w, data)
}

func (e Engine) Index(ctx context.Context, w http.ResponseWriter) {
	e.render(ctx, "views/index.tmpl", w, nil)
}

func (e Engine) render(
	ctx context.Context, path string, w http.ResponseWriter, data any,
) {
	log := utils.LoggerFromCtx(ctx).With(
		zap.Any("data", data),
		zap.String("path", path),
	)

	session := sessiondata.Read(ctx)
	viewData, err := makeViewData(session, data)
	if err != nil {
		log.Error("failed to make view data", zap.Error(err))
	}

	if err := e.Render(w, path, viewData); err != nil {
		log.Error("failed to render", zap.Error(err))
	}

}

func makeViewData(data ...any) (viewData, error) {
	var result viewData
	for _, d := range data {
		if err := mapstructure.Decode(d, &result); err != nil {
			return nil, err
		}
	}

	return result, nil
}
