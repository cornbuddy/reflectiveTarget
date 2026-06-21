package utils

import (
	"go.uber.org/zap"
)

// this function should make linter happier. [ResponseWriter.Write] returns int
// and error, and instead of checking for err != nil each time, just decorate
// writes with this function
func LogBadWrites(log *zap.Logger) func(int, error) {
	return func(n int, err error) {
		if err != nil {
			log.Error("failed to write response", zap.Error(err))
		}
	}
}
