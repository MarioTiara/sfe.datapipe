package logger

import "context"

type Field struct {
	Key   string
	Value interface{}
}

type Logger interface {
	Info(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, err error, fields ...Field)
	WithFields(fields ...Field) Logger
}
