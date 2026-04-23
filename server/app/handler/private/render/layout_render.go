package render

import (
	"context"
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
)

type layoutRender struct{}

func (layoutRender) Index(ctx context.Context, w http.ResponseWriter) {
	render(ctx, w, layout(sessiondata.Read(ctx), index()))
}

func (layoutRender) Login(
	ctx context.Context, w http.ResponseWriter, data LoginData,
) {
	render(ctx, w, layout(sessiondata.Read(ctx), login(data)))
}

func (layoutRender) Signup(
	ctx context.Context, w http.ResponseWriter, data SignupData,
) {
	render(ctx, w, layout(sessiondata.Read(ctx), signup(data)))
}

func (layoutRender) Targets(
	ctx context.Context, w http.ResponseWriter, data TargetsData,
) {
	render(ctx, w, layout(sessiondata.Read(ctx), targets(data)))
}

func (layoutRender) TargetForm(
	ctx context.Context, w http.ResponseWriter, data TargetFormData,
) {
	render(ctx, w, layout(sessiondata.Read(ctx), targetForm(data)))
}
