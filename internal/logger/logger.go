package logger

import (
	"go.uber.org/zap"
)

// type AppLogger struct {
// 	Logger *zap.SugaredLogger
// }

//	func NewAppLogger() *AppLogger {
//		logger, _ := zap.NewProduction()
//		defer logger.Sync()
//		sugar := logger.Sugar()
//		sugar.Infow("Логер создан")
//		return &AppLogger{Logger: sugar}
//	}
func NewAppLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Infow("Логер создан")
	return sugar, err
}
