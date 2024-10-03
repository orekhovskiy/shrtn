package logger

import "go.uber.org/zap"

type ZapLogger struct {
	*zap.Logger
}

func (z *ZapLogger) Info(msg string, fields ...zap.Field) {
	z.Logger.Info(msg, fields...)
}

func (z *ZapLogger) Error(msg string, fields ...zap.Field) {
	z.Logger.Error(msg, fields...)
}

func NewZapLogger() (*ZapLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{Logger: logger}, nil
}
