package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sort"
	"time"
)

// 配置信息
type Config struct {
	Program           string                 `json:"program" yaml:"program"`                     // 应用程序名称
	Level             AtomicLevel            `json:"level" yaml:"level"`                         // 记录的最低级别
	Development       bool                   `json:"development" yaml:"development"`             // 是否开发模式
	DisableCaller     bool                   `json:"disableCaller" yaml:"disableCaller"`         // 是否禁用获取调用者信息
	DisableStacktrace bool                   `json:"disableStacktrace" yaml:"disableStacktrace"` // 是否禁用获取堆栈跟踪
	DisableSpanLogger bool                   `json:"disableSpanLogger" yaml:"disableSpanLogger"` // 是否禁用Tracing的Span Logger
	SpanLevel         AtomicLevel            `json:"spanLevel" yaml:"spanLevel"`                 // Tracing Span Logger 记录的最低级别
	EncoderConfig     EncoderConfig          `json:"encoderConfig" yaml:"encoderConfig"`         // 日志编码器配置 推荐使用标准编码器配置
	Sampling          *SamplingConfig        `json:"sampling" yaml:"sampling"`                   // 日志采样配置 目的是限制住CPU/IO的负载 此处配置的值是每秒
	InitialFields     map[string]interface{} `json:"initialFields" yaml:"initialFields"`         // 记录器自带内置上下文字段

	StdoutEnabled  bool            `json:"stdoutEnabled" yaml:"stdoutEnabled"`   // 是否开启标准输出
	FileEnabled    bool            `json:"fileEnabled" yaml:"fileEnabled"`       // 是否开启文件输出
	FileCoreConfig *FileCoreConfig `json:"fileCoreConfig" yaml:"fileCoreConfig"` // 文件输出配置
}

// 采样配置
type SamplingConfig = zap.SamplingConfig

// 日志编码器配置
type EncoderConfig = zapcore.EncoderConfig

// 标准的编码器配置
func NewStandardEncoderConfig() EncoderConfig {
	return EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// 默认配置
func DefaultConfig() *Config {
	return ProductionConfig()
}

// 生产环境推荐配置
func ProductionConfig() *Config {

	cfg := &Config{
		Level:          NewAtomicLevelAt(InfoLevel),
		SpanLevel:      NewAtomicLevelAt(InfoLevel),
		Development:    false,
		EncoderConfig:  NewStandardEncoderConfig(),
		FileEnabled:    true,
		FileCoreConfig: StandardFileCoreConfig(),
	}

	return cfg
}

// 开发环境推荐配置
func DevelopmentConfig() *Config {

	cfg := &Config{
		Level:          NewAtomicLevelAt(DebugLevel),
		SpanLevel:      NewAtomicLevelAt(DebugLevel),
		Development:    true,
		EncoderConfig:  NewStandardEncoderConfig(),
		FileEnabled:    true,
		FileCoreConfig: StandardFileCoreConfig(),
	}

	return cfg
}

// 通过配置构建记录器
func (cfg Config) Build(opts ...Option) (*Logger, error) {

	var cores []zapcore.Core

	if cfg.FileEnabled {
		if cfg.FileCoreConfig != nil {
			fileCore, err := FileCore(cfg.FileCoreConfig, cfg.EncoderConfig, cfg.Level)
			if err != nil {
				return nil, err
			}
			cores = append(cores, fileCore)
		} else {
			fileCore, err := FileCore(StandardFileCoreConfig(), cfg.EncoderConfig, cfg.Level)
			if err != nil {
				return nil, err
			}
			cores = append(cores, fileCore)
		}
	}

	// 标准输出判断
	// 如果没有输出，则默认为标准输出
	if cfg.StdoutEnabled || len(cores) == 0 {
		cores = append(cores, StdCore(cfg.EncoderConfig, cfg.Level))
	}

	base := zap.New(zapcore.NewTee(cores...), cfg.buildOptions()...)

	logger := newLogger(base)

	if !cfg.DisableStacktrace {
		stackLevel := ErrorLevel
		if cfg.Development {
			stackLevel = WarnLevel
		}
		logger.addStack = stackLevel
	}

	if !cfg.DisableSpanLogger {
		logger.spanLevel = cfg.SpanLevel
	}

	if len(opts) > 0 {
		logger = logger.WithOptions(opts...)
	}

	return logger, nil
}

// 构建记录器的配置项
func (cfg Config) buildOptions() []zap.Option {
	var opts []zap.Option

	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(2))
	}

	if cfg.Sampling != nil {
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSampler(core, time.Second, int(cfg.Sampling.Initial), int(cfg.Sampling.Thereafter))
		}))
	}

	if len(cfg.InitialFields) > 0 {
		fs := make([]Field, 0, len(cfg.InitialFields))
		keys := make([]string, 0, len(cfg.InitialFields))
		for k := range cfg.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, zap.Any(k, cfg.InitialFields[k]))
		}
		opts = append(opts, zap.Fields(fs...))
	}

	if cfg.Program != "" {
		opts = append(opts, zap.Fields(program(cfg.Program)))
	}

	return opts
}
