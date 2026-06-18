package render

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/utils"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
)

type Render interface {
	Index(context.Context, http.ResponseWriter)
	Login(context.Context, http.ResponseWriter, LoginData)
	Signup(context.Context, http.ResponseWriter, SignupData)
	Targets(context.Context, http.ResponseWriter, TargetsData)
	TargetForm(context.Context, http.ResponseWriter, TargetFormData)
}

type TargetFormData struct {
	contracts.TargetForm
	valueobjects.ID
}

type TargetsData struct {
	aggregations.Targets
}

type SignupData struct {
	contracts.SignupForm
}

type LoginData struct {
	contracts.LoginForm
}

// renders view with layout (eg navbar, head, header, footer, etc)
var Layout Render = layoutRender{}

// renders view only
var View Render = viewRender{}

func render(ctx context.Context, w http.ResponseWriter, cmp templ.Component) {
	if err := cmp.Render(ctx, w); err != nil {
		log := log.Logger(ctx, zap.Error(err))
		log.Error("rendering failed")
		utils.HttpError(w, http.StatusInternalServerError)
	}
}
