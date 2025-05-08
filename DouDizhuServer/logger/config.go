package logger

import (
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// LogConfig 日志配置结构
type LogConfig struct {
	Debug   bool   `yaml:"debug"`    // 是否开启debug模式
	LogPath string `yaml:"log_path"` // 日志文件路径
}

// DefaultConfig 返回默认配置
func DefaultConfig() *LogConfig {
	return &LogConfig{
		Debug:   true,
		LogPath: "logs",
	}
}

// LoadConfig 从文件加载配置
func LoadConfig(path string) (*LogConfig, error) {
	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// 读取配置文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// 解析配置
	config := DefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// NewConfig 创建配置
func NewConfig(debug bool) zap.Config {
	if debug {
		cfg := zap.NewDevelopmentConfig()
		return cfg
	}

	cfg := zap.NewProductionConfig()
	return cfg
}
