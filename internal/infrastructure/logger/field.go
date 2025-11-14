package logger

import (
	"context"

	"github.com/mariotiara/sfe-data-pipe/internal/shared/logger"
	"github.com/mariotiara/sfe-data-pipe/internal/shared/processid"
	"go.uber.org/zap"
)

func toZapFields(fields ...logger.Field) []zap.Field {
	zFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zFields = append(zFields, zap.Any(f.Key, f.Value))
	}
	return zFields
}

func processidIDField(ctx context.Context) logger.Field {
	id := processid.FromContext(ctx)
	if id == "" {
		return logger.Field{Key: "process_id", Value: "N/A"}
	}
	return logger.Field{Key: "process_id", Value: id}
}
