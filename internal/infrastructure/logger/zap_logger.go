package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zapLogger struct {
	base *zap.Logger
}

func NewZapLogger() logger.Logger {
	// 1. Ensure logs folder exists
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.MkdirAll("logs", 0755)
	}

	// 2. Setup lumberjack for rolling file
	filename := fmt.Sprintf("logs/app-%s.log", time.Now().Format("2006-01-02"))

	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10, // MB
		MaxBackups: 7,
		MaxAge:     30, // days
		Compress:   false,
	}

	// 3. Encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)
	level := zapcore.InfoLevel

	// 4. MultiWriteSyncer: file + console
	fileSyncer := zapcore.AddSync(lumberJackLogger)
	consoleSyncer := zapcore.AddSync(os.Stdout)
	multiSyncer := zapcore.NewMultiWriteSyncer(fileSyncer, consoleSyncer)

	// 5. Core
	core := zapcore.NewCore(encoder, multiSyncer, level)

	// 6. Logger with caller info
	l := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(1))

	return &zapLogger{base: l}
}

// Info logs info-level messages
func (l *zapLogger) Info(ctx context.Context, msg string, fields ...logger.Field) {
	fields = append(fields, processidIDField(ctx))
	zapFields := toZapFields(fields...)
	l.base.Info(msg, zapFields...)
}

// Error logs error-level messages
func (l *zapLogger) Error(ctx context.Context, msg string, err error, fields ...logger.Field) {
	fields = append(fields, processidIDField(ctx))
	fields = append(fields, logger.Field{Key: "error", Value: err.Error()})
	zapFields := toZapFields(fields...)
	l.base.Error(msg, zapFields...)
}

// WithFields adds fields to the logger
func (l *zapLogger) WithFields(fields ...logger.Field) logger.Logger {
	return &zapLogger{base: l.base.With(toZapFields(fields...)...)}
}
