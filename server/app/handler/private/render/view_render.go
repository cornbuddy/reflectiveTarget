package render

import (
	"context"
	"net/http"
)

type viewRender struct{}

func (viewRender) Index(ctx context.Context, w http.ResponseWriter) error {
	return index().Render(ctx, w)
}

func (viewRender) Login(
	ctx context.Context, w http.ResponseWriter, data LoginData,
) error {
	return login(data).Render(ctx, w)
}

func (viewRender) Signup(
	ctx context.Context, w http.ResponseWriter, data SignupData,
) error {
	return signup(data).Render(ctx, w)
}

func (viewRender) Targets(
	ctx context.Context, w http.ResponseWriter, data TargetsData,
) error {
	return targets(data).Render(ctx, w)
}

func (viewRender) TargetForm(
	ctx context.Context, w http.ResponseWriter, data TargetFormData,
) error {
	return targetForm(data).Render(ctx, w)
}
