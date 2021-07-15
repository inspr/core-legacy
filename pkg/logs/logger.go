package logs

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var alevel zap.AtomicLevel

// Logger returns the singleton logger and it's atomic level
func Logger(options ...zap.Option) (*zap.Logger, *zap.AtomicLevel) {
	if logger == nil {
		var level zapcore.Level
		switch os.Getenv("LOG_LEVEL") {
		case "warn", "warning":
			level = zap.WarnLevel
		case "debug":
			level = zap.DebugLevel
		case "error":
			level = zap.ErrorLevel
		}
		alevel = zap.NewAtomicLevelAt(level)
		enablerFunc := func(l zapcore.Level) bool {
			return l >= alevel.Level()
		}
		core := zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()), os.Stdout, zap.LevelEnablerFunc(enablerFunc))
		logger = zap.New(core)
	}

	return logger.WithOptions(options...), &alevel
}
