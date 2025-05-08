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
func Info(args ...interface{}) {
	sugar.Info(args...)
}

// Infof 输出info级别日志（格式化）
func Infof(format string, args ...interface{}) {
	sugar.Infof(format, args...)
}

// Error 输出error级别日志
func Error(args ...interface{}) {
	sugar.Error(args...)
}

// Errorf 输出error级别日志（格式化）
func Errorf(format string, args ...interface{}) {
	sugar.Errorf(format, args...)
}

// Debug 输出debug级别日志
func Debug(args ...interface{}) {
	sugar.Debug(args...)
}

// Debugf 输出debug级别日志（格式化）
func Debugf(format string, args ...interface{}) {
	sugar.Debugf(format, args...)
}

// Warn 输出warn级别日志
func Warn(args ...interface{}) {
	sugar.Warn(args...)
}

// Warnf 输出warn级别日志（格式化）
func Warnf(format string, args ...interface{}) {
	sugar.Warnf(format, args...)
}

// Fatal 输出fatal级别日志并退出程序
func Fatal(args ...interface{}) {
	sugar.Fatal(args...)
}

// Fatalf 输出fatal级别日志并退出程序（格式化）
func Fatalf(format string, args ...interface{}) {
	sugar.Fatalf(format, args...)
}

// Panic 输出panic级别日志
func Panic(args ...interface{}) {
	sugar.Panic(args...)
}

// Panicf 输出panic级别日志（格式化）
func Panicf(format string, args ...interface{}) {
	sugar.Panicf(format, args...)
}

// DPanic 输出dpanic级别日志
func DPanic(args ...interface{}) {
	sugar.DPanic(args...)
}

// DPanicf 输出dpanic级别日志（格式化）
func DPanicf(format string, args ...interface{}) {
	sugar.DPanicf(format, args...)
}

// Sync 同步日志
func Sync() error {
	return sugar.Sync()
}
