package logger

import "go.uber.org/zap/zapcore"

type Encoder = zapcore.Encoder

type ObjectEncoder = zapcore.ObjectEncoder

func NewJSONEncoder(cfg EncoderConfig) Encoder {
	return zapcore.NewJSONEncoder(cfg)
}
