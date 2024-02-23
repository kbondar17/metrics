package logger

import (
	"go.uber.org/zap"
)

func New() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Infow("Логер создан")
	return sugar, err
}
