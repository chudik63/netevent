package logger

import (
	"context"

	"go.uber.org/zap"
)

const (
	ServiceName = "auth"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

type logger struct {
	serviceName string
	logger      *zap.Logger
}

func (l logger) Info(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("ServiceName", ServiceName))
	l.logger.Info(msg, fields...)
}

func (l logger) Error(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("ServiceName", ServiceName))
	l.logger.Error(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	fields = append(fields, zap.String("ServiceName", "auth"))
	zapLogger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	fields = append(fields, zap.String("ServiceName", "auth"))
	zapLogger.Error(msg, fields...)
}

func New(serviceName string) Logger {
	zapLogger, _ := zap.NewDevelopment()

	defer zapLogger.Sync()
	return &logger{
		serviceName: serviceName,
		logger:      zapLogger,
	}
}

func CtxGetLogger(ctx context.Context) Logger {
	return ctx.Value(ServiceName).(Logger)
}
