package render

import (
	"github.com/abiosoft/mold"

	"github.com/cornbuddy/reflectiveTarget/server/app/forms"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
)

type TargetsData struct {
	aggregations.Targets
}

type SignupData struct {
	forms.SignupForm
}

type LoginData struct {
	forms.LoginForm
}

type Engine struct {
	mold.Engine
}

type viewData map[string]any
