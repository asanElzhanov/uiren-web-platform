package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger(logLevel string) error {
	level, err := parseLogLevel(logLevel)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	zapCfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Encoding:    "json",
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:    "timestamp",
			LevelKey:   "level",
			MessageKey: "message",
			CallerKey:  "caller",

			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeTime:   zapcore.RFC3339TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	Log, err = zapCfg.Build()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	return nil
}

func parseLogLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("unknown log level: %s", level)
	}
}

func Info(msg string, args ...interface{}) {
	Log.Info(msg, convertArgsToZapFields(args...)...)
}

func Error(msg string, args ...interface{}) {
	Log.Error(msg, convertArgsToZapFields(args...)...)
}

func Warn(msg string, args ...interface{}) {
	Log.Warn(msg, convertArgsToZapFields(args...)...)
}

func Fatal(msg string, args ...interface{}) {
	Log.Fatal(msg, convertArgsToZapFields(args...)...)
}

func convertArgsToZapFields(args ...interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(args))

	for _, arg := range args {
		switch v := arg.(type) {
		case zap.Field:
			fields = append(fields, v)
		case string:
			fields = append(fields, zap.String("message", v))
		case int:
			fields = append(fields, zap.Int("value", v))
		case bool:
			fields = append(fields, zap.Bool("flag", v))
		case error:
			fields = append(fields, zap.Error(v))
		/*case map[string]interface{}:
		fields = append(fields, MapToZapFields(v)...)*/
		default:
			fields = append(fields, zap.Any("data", v))
		}
	}

	return fields
}
