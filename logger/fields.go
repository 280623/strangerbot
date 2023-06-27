package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zapcore.Field

func Any(key string, value interface{}) Field {
	return zap.Any(key, value)
}

func Binary(key string, value []byte) Field {
	return zap.Binary(key, value)
}

func Bool(key string, value bool) Field {
	return zap.Bool(key, value)
}

func ByteString(key string, value []byte) Field {
	return zap.ByteString(key, value)
}

func Complex128(key string, value complex128) Field {
	return zap.Complex128(key, value)
}

func Complex64(key string, value complex64) Field {
	return zap.Complex64(key, value)
}

func Duration(key string, value time.Duration) Field {
	return zap.Duration(key, value)
}

func Error(err error) Field {
	return zap.Error(err)
}

func Errors(key string, errs []error) Field {
	return zap.Errors(key, errs)
}

func NamedError(key string, err error) Field {
	return zap.NamedError(key, err)
}

func Float32(key string, value float32) Field {
	return zap.Float32(key, value)
}

func Float64(key string, value float64) Field {
	return zap.Float64(key, value)
}

func Int(key string, value int) Field {
	return zap.Int(key, value)
}

func Int8(key string, value int8) Field {
	return zap.Int8(key, value)
}

func Int16(key string, value int16) Field {
	return zap.Int16(key, value)
}

func Int32(key string, value int32) Field {
	return zap.Int32(key, value)
}

func Int64(key string, value int64) Field {
	return zap.Int64(key, value)
}

func Uint(key string, value uint) Field {
	return zap.Uint(key, value)
}

func Uint8(key string, value uint8) Field {
	return zap.Uint8(key, value)
}

func Uint16(key string, value uint16) Field {
	return zap.Uint16(key, value)
}

func Uint32(key string, value uint32) Field {
	return zap.Uint32(key, value)
}

func Uint64(key string, value uint64) Field {
	return zap.Uint64(key, value)
}

func String(key string, value string) Field {
	return zap.String(key, value)
}

func Stringer(key string, value fmt.Stringer) Field {
	return zap.Stringer(key, value)
}

func Strings(key string, value []string) Field {
	return zap.Strings(key, value)
}

func Time(key string, value time.Time) Field {
	return zap.Time(key, value)
}

func Uintptr(key string, value uintptr) Field {
	return zap.Uintptr(key, value)
}

func Object(key string, val zapcore.ObjectMarshaler) Field {
	return zap.Object(key, val)
}

func Array(key string, val zapcore.ArrayMarshaler) Field {
	return zap.Array(key, val)
}

func Stack(key string) Field {
	return zap.Stack(key)
}

func Reflect(key string, val interface{}) Field {
	return zap.Reflect(key, val)
}
