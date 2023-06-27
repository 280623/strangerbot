package logger

import (
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap"
)

type Level = zapcore.Level

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel = zapcore.DebugLevel
	// InfoLevel is the default logging priority.
	InfoLevel = zapcore.InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel = zapcore.WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel = zapcore.ErrorLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zapcore.FatalLevel
)

// AtomicLevel
type AtomicLevel = zap.AtomicLevel

type LevelEnabler = zapcore.LevelEnabler

func NewAtomicLevelAt(l Level) AtomicLevel {
	return zap.NewAtomicLevelAt(l)
}
