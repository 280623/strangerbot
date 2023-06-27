package logger

import (
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"go.uber.org/zap"
)

// 文件输出配置
type FileCoreConfig struct {
	Filename      string `json:"filename" yaml:"filename"`           // 保存的文件路径
	DisableRotate bool   `json:"disableRotate" yaml:"disableRotate"` // 是否禁用轮转策略 另外使用系统脚本去控制轮转
	MaxSize       int    `json:"maxSize" yaml:"maxSize"`             // 单个日志文件最大存储 单位MB
	MaxBackups    int    `json:"maxBackups" yaml:"maxBackups"`       // 最多保存的文件数
	MaxAge        int    `json:"maxAge" yaml:"maxAge"`               // 最多保存的天数
}

// 默认文件输出配置
func StandardFileCoreConfig() *FileCoreConfig {
	return &FileCoreConfig{Filename: "./logs/app.log", MaxSize: 50, MaxBackups: 30, MaxAge: 30}
}

// 文件输出
func FileCore(fileCoreConfig *FileCoreConfig, config zapcore.EncoderConfig, enab zapcore.LevelEnabler) (zapcore.Core, error) {

	encoder := zapcore.NewJSONEncoder(config)

	if !fileCoreConfig.DisableRotate {

		writerSync := zapcore.AddSync(&lumberjack.Logger{
			Filename:   fileCoreConfig.Filename,
			MaxSize:    fileCoreConfig.MaxSize,
			MaxBackups: fileCoreConfig.MaxBackups,
			MaxAge:     fileCoreConfig.MaxAge,
			LocalTime:  true,
		})

		core := zapcore.NewCore(encoder, writerSync, enab)

		return core, nil
	}

	writerSync, _, err := zap.Open(fileCoreConfig.Filename)
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(encoder, writerSync, enab)

	return core, nil
}
