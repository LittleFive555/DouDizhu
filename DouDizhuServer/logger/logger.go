package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

var (
	sugar *zap.SugaredLogger
)

// InitLoggerWithConfig 使用配置文件初始化logger
func InitLoggerWithConfig(configPath string) error {
	// 加载配置
	config, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	// 创建zap配置
	cfg := NewConfig(config.Debug)

	// 设置日志文件
	if err := os.MkdirAll(config.LogPath, 0755); err != nil {
		return err
	}

	cfg.OutputPaths = []string{
		filepath.Join(config.LogPath, "app.log"),
		"stdout",
	}
	cfg.ErrorOutputPaths = []string{
		filepath.Join(config.LogPath, "error.log"),
		"stdout",
	}

	// 构建logger，跳过我们的包装层
	logger, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}
	sugar = logger.Sugar()
	return nil
}

// With 创建带有字段的logger
func With(args ...interface{}) *zap.SugaredLogger {
	return sugar.With(args...)
}

// WithField 创建带有单个字段的logger
func WithField(key string, value interface{}) *zap.SugaredLogger {
	return sugar.With(key, value)
}

// WithFields 创建带有多个字段的logger
func WithFields(fields map[string]interface{}) *zap.SugaredLogger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return sugar.With(args...)
}

// Info 输出info级别日志
func Info(msg string) {
	sugar.Info(msg)
}

// Infof 输出info级别日志（格式化）
func InfoFormat(format string, args ...interface{}) {
	sugar.Infof(format, args...)
}

func InfoWith(msg string, fields ...interface{}) {
	sugar.Infow(msg, fields...)
}

// Error 输出error级别日志
func Error(msg string) {
	sugar.Error(msg)
}

// Errorf 输出error级别日志（格式化）
func ErrorFormat(format string, args ...interface{}) {
	sugar.Errorf(format, args...)
}

func ErrorWith(msg string, fields ...interface{}) {
	sugar.Errorw(msg, fields...)
}

// Debug 输出debug级别日志
func Debug(msg string) {
	sugar.Debug(msg)
}

// Debugf 输出debug级别日志（格式化）
func DebugFormat(format string, args ...interface{}) {
	sugar.Debugf(format, args...)
}

func DebugWith(msg string, fields ...interface{}) {
	sugar.Debugw(msg, fields...)
}

// Warn 输出warn级别日志
func Warn(msg string) {
	sugar.Warn(msg)
}

// Warnf 输出warn级别日志（格式化）
func WarnFormat(format string, args ...interface{}) {
	sugar.Warnf(format, args...)
}

func WarnWith(msg string, fields ...interface{}) {
	sugar.Warnw(msg, fields...)
}

// Fatal 输出fatal级别日志并退出程序
func Fatal(msg string) {
	sugar.Fatal(msg)
}

// Fatalf 输出fatal级别日志并退出程序（格式化）
func FatalFormat(format string, args ...interface{}) {
	sugar.Fatalf(format, args...)
}

func FatalWith(msg string, fields ...interface{}) {
	sugar.Fatalw(msg, fields...)
}

// Panic 输出panic级别日志
func Panic(msg string) {
	sugar.Panic(msg)
}

// Panicf 输出panic级别日志（格式化）
func PanicFormat(format string, args ...interface{}) {
	sugar.Panicf(format, args...)
}

func PanicWith(msg string, fields ...interface{}) {
	sugar.Panicw(msg, fields...)
}

// DPanic 输出dpanic级别日志
func DPanic(msg string) {
	sugar.DPanic(msg)
}

// DPanicf 输出dpanic级别日志（格式化）
func DPanicFormat(format string, args ...interface{}) {
	sugar.DPanicf(format, args...)
}

func DPanicWith(msg string, fields ...interface{}) {
	sugar.DPanicw(msg, fields...)
}

// Sync 同步日志
func Sync() error {
	return sugar.Sync()
}
