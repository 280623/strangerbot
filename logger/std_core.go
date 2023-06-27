package logger

import (
	"os"

	"go.uber.org/zap/zapcore"
)

func StdCore(config zapcore.EncoderConfig, enab zapcore.LevelEnabler) zapcore.Core {

	encoder := zapcore.NewJSONEncoder(config)
	writer := zapcore.Lock(os.Stdout)
	core := zapcore.NewCore(encoder, writer, enab)

	return core
}
