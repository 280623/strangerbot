package logger

import "go.uber.org/zap/zapcore"

type Core = zapcore.Core
type WriteSyncer = zapcore.WriteSyncer

// 创建新的核心
// level enabler 是给core写入的时候做最低级别日志的判断
// span level 不在底层core 在logger层，顾没有传入
func NewCore(enc Encoder, ws WriteSyncer, enab LevelEnabler) Core {
	core := zapcore.NewCore(enc, ws, enab)
	return core
}

// 创建多个底层核心
// 日志记录会复制到更多的核心内
func NewTee(cores ...Core) Core {
	return zapcore.NewTee(cores...)
}
