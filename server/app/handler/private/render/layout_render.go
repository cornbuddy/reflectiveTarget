package render

import (
	"context"
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
)

type layoutRender struct{}

func (layoutRender) Index(ctx context.Context, w http.ResponseWriter) error {
	return layout(sessiondata.Read(ctx), index()).Render(ctx, w)
}

func (layoutRender) Login(
	ctx context.Context, w http.ResponseWriter, data LoginData,
) error {
	return layout(sessiondata.Read(ctx), login(data)).Render(ctx, w)
}

func (layoutRender) Signup(
	ctx context.Context, w http.ResponseWriter, data SignupData,
) error {
	return layout(sessiondata.Read(ctx), signup(data)).Render(ctx, w)
}

func (layoutRender) Targets(
	ctx context.Context, w http.ResponseWriter, data TargetsData,
) error {
	return layout(sessiondata.Read(ctx), targets(data)).Render(ctx, w)
}

func (layoutRender) TargetForm(
	ctx context.Context, w http.ResponseWriter, data TargetFormData,
) error {
	return layout(sessiondata.Read(ctx), targetForm(data)).Render(ctx, w)
}
