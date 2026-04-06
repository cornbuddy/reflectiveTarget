package render

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/app/utils"
)

func (e engine) Targets(
	ctx context.Context, w http.ResponseWriter, data TargetsData,
) {
	e.render(ctx, "views/targets.tmpl", w, data)
}

func (e engine) Signup(
	ctx context.Context, w http.ResponseWriter, data SignupData,
) {
	e.render(ctx, "views/signup.tmpl", w, data)
}

func (e engine) Login(
	ctx context.Context, w http.ResponseWriter, data LoginData,
) {
	e.render(ctx, "views/login.tmpl", w, data)
}

func (e engine) Index(ctx context.Context, w http.ResponseWriter) {
	e.render(ctx, "views/index.tmpl", w, nil)
}

func (e engine) render(
	ctx context.Context, path string, w http.ResponseWriter, data any,
) {
	log := utils.LoggerFromCtx(ctx).With(
		zap.Any("data", data),
		zap.String("path", path),
	)

	session := sessiondata.Read(ctx)
	viewData, err := makeViewData(data, session)
	log = log.With(zap.Any("view-data", viewData))
	log.Debug("view data")
	if err != nil {
		log.Error("failed to make view data", zap.Error(err))
	}

	if err := e.Render(w, path, *viewData); err != nil {
		log.Error("failed to render", zap.Error(err))
	}

}

func makeViewData(datas ...any) (*viewData, error) {
	var result viewData
	for _, data := range datas {
		switch v := data.(type) {
		case sessiondata.SessionData:
			result.SessionData = v
		case TargetsData:
			result.TargetsData = v
		case SignupData:
			result.SignupData = v
		case LoginData:
			result.LoginData = v
		case nil:
			// do nothing, this is expected
		default:
			return nil, fmt.Errorf("unknown data %v", data)
		}
	}

	return &result, nil
}
