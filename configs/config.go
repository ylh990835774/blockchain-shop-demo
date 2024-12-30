package configs

import (
	"github.com/spf13/viper"
)

// Config 是应用程序的配置
type Config struct {
	Server ServerConfig `yaml:"server"`
	MySQL  MySQLConfig  `yaml:"mysql"`
	JWT    JWTConfig    `yaml:"jwt"`
	Log    LogConfig    `yaml:"log"`
}

// ServerConfig 是服务器配置
type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

// MySQLConfig 是MySQL数据库配置
type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// JWTConfig 是JWT配置
type JWTConfig struct {
	SecretKey           string `yaml:"secretKey"`
	Issuer              string `yaml:"issuer"`
	ExpireDurationHours int    `yaml:"expireDurationHours"`
}

// LogConfig 是日志配置
type LogConfig struct {
	Level      string `yaml:"level"`      // 日志级别: debug, info, warn, error, fatal
	Format     string `yaml:"format"`     // 日志格式: json, console
	Filename   string `yaml:"filename"`   // 日志文件路径
	MaxSize    int    `yaml:"maxSize"`    // 每个日志文件最大尺寸，单位MB
	MaxBackups int    `yaml:"maxBackups"` // 保留的旧日志文件最大数量
	MaxAge     int    `yaml:"maxAge"`     // 保留的旧日志文件最大天数
	Compress   bool   `yaml:"compress"`   // 是否压缩旧日志文件
	Console    bool   `yaml:"console"`    // 是否同时输出到控制台
}

// Load 加载配置文件
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
