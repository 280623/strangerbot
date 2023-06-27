package logger

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ILogger interface {
	Named(name string) ILogger
	WithOptions(opts ...Option) ILogger
	With(fields ...Field) ILogger
	Program(name string) ILogger
	Span(span opentracing.Span) ILogger
	SpanFromContext(ctx context.Context) ILogger
	ShortLogger
	Core() Core
	Sync()
}

type ShortLogger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}

//
//var _ ILogger = &Logger{}

type Logger struct {
	base      *zap.Logger
	addStack  zapcore.LevelEnabler
	span      opentracing.Span
	spanLevel zapcore.LevelEnabler
}

// 开发模式记录器
func NewDevelopment(options ...Option) (*Logger, error) {
	return DevelopmentConfig().Build(options...)
}

// 生产模式记录器
func NewProduction(options ...Option) (*Logger, error) {
	return ProductionConfig().Build(options...)
}

// 通过zap实例来生成记录器
func newLogger(base *zap.Logger) *Logger {

	logger := &Logger{
		base:      base,
		addStack:  zapcore.FatalLevel + 1,
		spanLevel: zapcore.FatalLevel + 1,
	}

	return logger
}

// Named 增加记录器的名称路径 如已存在则使用追加路径
// 默认情况下，记录器未命名
func (l *Logger) Named(name string) *Logger {

	if name == "" {
		return l
	}

	cp := l.clone()
	cp.base = l.base.Named(name)

	return cp
}

// WithOptions克隆当前的Logger，应用提供的Options
// 返回生成的Logger 并发安全
func (l *Logger) WithOptions(opts ...Option) *Logger {

	cp := l.clone()
	for _, opt := range opts {
		opt.apply(cp)
	}

	return cp

}

// 创建子记录器并向其添加结构化上下文字段
// 子记录器增加字段不影响父记录，反之亦然
func (l *Logger) With(fields ...Field) *Logger {

	if len(fields) == 0 {
		return l
	}

	cp := l.clone()
	cp.base = cp.base.With(fields...)

	return cp
}

// 为记录器添加程序名称说明
// 如果传入名称不为空，则创建子记录器含程序名称
func (l *Logger) Program(name string) *Logger {

	if name == "" {
		return l
	}

	cp := l.clone()
	cp.base = cp.base.With(program(name))

	return cp
}

// 为记录器关联Tracing Span对象
// 如果span不存在则返回原记录器，如果span存在则返回新创建的子记录器
func (l *Logger) Span(span opentracing.Span) *Logger {

	if span != nil {
		cp := l.clone()
		cp.span = span
		return cp
	}

	return l
}

// 记录器通过请求上下文关联Tracing Span对象
// 如果span不存在则返回原记录器，如果span存在则返回新创建的子记录器
func (l *Logger) SpanFromContext(ctx context.Context) *Logger {

	if span := opentracing.SpanFromContext(ctx); span != nil {
		return l.Span(span)
	}

	return l
}

// 打印DEBUG级别日志 包含fields所有字段
// 以及记录器已累积的字段
func (l *Logger) Debug(msg string, fields ...Field) {
	l.log(DebugLevel, "debug", msg, fields)
}

// 打印INFO级别日志 包含fields所有字段
// 以及记录器已累积的字段
func (l *Logger) Info(msg string, fields ...Field) {
	l.log(InfoLevel, "info", msg, fields)
}

// 打印WARN级别日志 包含fields所有字段
// 以及记录器已累积的字段
func (l *Logger) Warn(msg string, fields ...Field) {
	l.log(WarnLevel, "warn", msg, fields)
}

// 打印ERROR级别日志 包含fields所有字段
// 以及记录器已累积的字段
func (l *Logger) Error(msg string, fields ...Field) {
	l.log(ErrorLevel, "error", msg, fields)
}

// 打印FATAL级别日志 包含fields所有字段
// 以及记录器已累积的字段
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.log(FatalLevel, "fatal", msg, fields)
}

// 返回记录器的底层zapcore.Core
func (l *Logger) Core() Core {
	return l.base.Core()
}

// 刷新已缓冲的日志
// 在应用程序退出之前应注意调用Sync
func (l *Logger) Sync() {
	l.base.Sync()
}

// clone logger
func (log *Logger) clone() *Logger {
	copy := *log
	return &copy
}

// 將日志打印到zap核心
func (l *Logger) log(zapLevel zapcore.Level, level string, msg string, fields []Field) {

	stackEnabled := false
	stack := ""
	if l.addStack.Enabled(zapLevel) {
		stackEnabled = true
		stack = zap.Stack("stacktrace").String
	}

	if ce := l.base.Check(zapLevel, msg); ce != nil {
		if stackEnabled {
			ce.Entry.Stack = stack
		}
		ce.Write(fields...)
	}

	l.logToSpan(zapLevel, level, msg, fields, stackEnabled, stack)

}

// 将日志打印到span logs
func (l *Logger) logToSpan(zapLevel zapcore.Level, level string, msg string, fields []Field, stackEnabled bool, stack string) {

	// 判断span是否为nil
	if l.span == nil {
		return
	}

	// 判断level是否允许打印
	if !l.spanLevel.Enabled(zapLevel) {
		return
	}

	// 将msg level fields转为span log field
	fa := fieldAdapter(make([]log.Field, 0, 2+len(fields)))
	fa = append(fa, log.String("event", msg))
	fa = append(fa, log.String("level", level))

	for _, field := range fields {
		field.AddTo(&fa)
	}

	if stackEnabled {
		fa = append(fa, log.String("stacktrace", stack))
	}

	// log to span
	l.span.LogFields(fa...)
}

// Program 字段
func program(program string) Field {
	return Field{Key: "program", Type: zapcore.StringType, String: program}
}
