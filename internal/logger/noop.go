package logger

import "go.uber.org/zap"

type NoopLogger struct {
	*zap.Logger
}

func (n *NoopLogger) Info(msg string, fields ...zap.Field) {}

func (n *NoopLogger) Error(msg string, fields ...zap.Field) {}
