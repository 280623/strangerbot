package logger

type EmptyShortLogger struct{}

func NewEmptyShortLogger() ShortLogger                        { return &EmptyShortLogger{} }
func (p *EmptyShortLogger) Debug(msg string, fields ...Field) {}
func (p *EmptyShortLogger) Info(msg string, fields ...Field)  {}
func (p *EmptyShortLogger) Warn(msg string, fields ...Field)  {}
func (p *EmptyShortLogger) Error(msg string, fields ...Field) {}
func (p *EmptyShortLogger) Fatal(msg string, fields ...Field) {}
