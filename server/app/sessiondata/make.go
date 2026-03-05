package sessiondata

import (
	"context"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/utils"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type SessionData struct {
	valueobjects.ID
	Username        string
	IsAuthetnicated bool
	Logger          *zap.Logger
}

func Make(ctx context.Context) SessionData {
	log := utils.LoggerFromCtx(ctx)

	return SessionData{
		Logger: log,
	}
}
