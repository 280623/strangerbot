package logger

var _defaultLogger *Logger

func GetLogger() *Logger {
	return _defaultLogger
}

func DeaultLoggerSync() {
	_defaultLogger.Sync()
}

func InitDefaultLogger(isDebug bool) error {

	cfg := ProductionConfig()
	if isDebug {
		cfg = DevelopmentConfig()
	}

	// 生成Logger
	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	_defaultLogger = logger

	return nil
}
