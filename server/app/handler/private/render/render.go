package render

import (
	"context"
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type Render interface {
	Index(context.Context, http.ResponseWriter) error
	Login(context.Context, http.ResponseWriter, LoginData) error
	Signup(context.Context, http.ResponseWriter, SignupData) error
	Targets(context.Context, http.ResponseWriter, TargetsData) error
	TargetForm(context.Context, http.ResponseWriter, TargetFormData) error
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
