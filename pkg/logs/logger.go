package logs

import (
	"log"
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
		err := alevel.UnmarshalText([]byte(os.Getenv("LOG_LEVEL")))
		if err != nil {
			log.Println("log level either not set or invalid, defaulting to info")
		}
		alevel = zap.NewAtomicLevelAt(level)
		enablerFunc := func(l zapcore.Level) bool {
			return l >= alevel.Level()
		}
		config := zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.RFC3339TimeEncoder
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
		core := zapcore.NewCore(zapcore.NewConsoleEncoder(config), os.Stdout, zap.LevelEnablerFunc(enablerFunc))
		logger = zap.New(core)
	}

	return logger.WithOptions(options...), &alevel
}
