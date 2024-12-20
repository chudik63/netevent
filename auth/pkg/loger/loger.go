package loger

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

type loger struct {
	serviceName string
	loger       *zap.Logger
}

func (l loger) Info(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("ServiceName", ServiceName))
	l.loger.Info(msg, fields...)
}

func (l loger) Error(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("ServiceName", ServiceName))
	l.loger.Error(msg, fields...)
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
	// zapLogger, _ := zap.NewProduction()
	zapLogger, _ := zap.NewDevelopment()

	defer zapLogger.Sync()
	return &loger{
		serviceName: serviceName,
		loger:       zapLogger,
	}
}

func CtxGetLogger(ctx context.Context) Logger {
	return ctx.Value(ServiceName).(Logger)
}
