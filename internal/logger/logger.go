package logger

import (
	"analytic-service/config"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	teamKey    = "Team21"
	serviceKey = "Analytic"
)

type Logger struct {
	*zap.Logger
}

func New(env *config.Config) (*Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoder := zapcore.NewJSONEncoder(cfg.EncoderConfig)
	logLevel, err := zapcore.ParseLevel(env.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("NewLogger: unexpected log level: %s", env.LogLevel)
	}

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), logLevel),
	)

	return &Logger{
		Logger: zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).With(
			zap.String("team", teamKey),
			zap.String("service", serviceKey),
		),
	}, nil
}

func Mock(testlogger *zap.Logger) *Logger {
	if testlogger == nil {
		return &Logger{
			Logger: zap.NewNop(),
		}
	}
	return &Logger{Logger: testlogger}
}
