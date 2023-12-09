package logger

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/adepte-myao/lamoda-test-2023/configs"
)

var (
	ErrUnknownLogLevel = errors.New("unknown log level")
)

func NewZap(config configs.AppConfig) (Logger, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig = encoderConfig
	logConfig.DisableStacktrace = !config.Logger.StackTraceEnabled

	level, ok := levelMapping[config.Logger.Level]
	if !ok {
		return nil, ErrUnknownLogLevel
	}
	logConfig.Level.SetLevel(level)

	coreLogger, err := logConfig.Build()
	if err != nil {
		return nil, err
	}

	return coreLogger.Sugar(), nil
}

var (
	levelMapping = map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"fatal": zapcore.FatalLevel,
	}
)
