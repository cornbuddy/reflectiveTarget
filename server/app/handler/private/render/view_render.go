package render

import (
	"context"
	"net/http"
)

type viewRender struct{}

func (viewRender) Index(ctx context.Context, w http.ResponseWriter) {
	render(ctx, w, index())
}

func (viewRender) Login(
	ctx context.Context, w http.ResponseWriter, data LoginData,
) {
	render(ctx, w, login(data))
}

func (viewRender) Signup(
	ctx context.Context, w http.ResponseWriter, data SignupData,
) {
	render(ctx, w, signup(data))
}

func (viewRender) Targets(
	ctx context.Context, w http.ResponseWriter, data TargetsData,
) {
	render(ctx, w, targets(data))
}

func (viewRender) TargetForm(
	ctx context.Context, w http.ResponseWriter, data TargetFormData,
) {
	render(ctx, w, targetForm(data))
}
