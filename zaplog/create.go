package zaplog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogName string

const (
	LogNameLogger   LogName = "logger"
	LogNameEcho     LogName = "echo"
	LogNameReductor LogName = "reductor"
	LogNameTrue     LogName = "true"
)

func isValidLogName(s string) bool {
	switch LogName(s) {
	case LogNameLogger, LogNameEcho, LogNameReductor, LogNameTrue:
		return true
	default:
		return false
	}
}

var encoderConfig = zapcore.EncoderConfig{
	TimeKey:        "ts",
	LevelKey:       "lvl",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

func createLogger(output []string, debug bool) (*zap.Logger, error) {
	level := zap.InfoLevel
	if debug {
		level = zap.DebugLevel
	}
	if len(output) == 0 {
		// Default to stdout when no outputs are provided.
		output = []string{"stdout"}
	}
	config := zap.Config{
		Level:         zap.NewAtomicLevelAt(level),
		Development:   debug,
		DisableCaller: false,
		// Only include stacktraces in debug mode to reduce noise in production logs.
		DisableStacktrace: !debug,
		// Human-friendly in debug, structured in production.
		Encoding: func() string {
			if debug {
				return "console"
			}
			return "json"
		}(),
		// Enable sampling in production to limit log volume.
		Sampling: func() *zap.SamplingConfig {
			if debug {
				return nil
			}
			return &zap.SamplingConfig{Initial: 100, Thereafter: 100}
		}(),
		EncoderConfig:    encoderConfig,
		OutputPaths:      output,
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
