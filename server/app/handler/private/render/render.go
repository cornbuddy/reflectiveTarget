package render

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/utils"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
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
		log := utils.LoggerFromCtx(ctx)
		utils.InternalServerError(log, w, "failed to render response", err)
	}
}
