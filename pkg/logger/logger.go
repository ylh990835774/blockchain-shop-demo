package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *zap.Logger

// Config 日志配置
type Config struct {
	Level      string `yaml:"level" json:"level"`             // 日志级别: debug, info, warn, error, fatal
	Format     string `yaml:"format" json:"format"`           // 日志格式: json, console
	Filename   string `yaml:"filename" json:"filename"`       // 日志文件路径
	MaxSize    int    `yaml:"max_size" json:"max_size"`       // 每个日志文件最大尺寸，单位MB
	MaxBackups int    `yaml:"max_backups" json:"max_backups"` // 保留的旧日志文件最大数量
	MaxAge     int    `yaml:"max_age" json:"max_age"`         // 保留的旧日志文件最大天数
	Compress   bool   `yaml:"compress" json:"compress"`       // 是否压缩旧日志文件
	Console    bool   `yaml:"console" json:"console"`         // 是否同时输出到控制台
}

// Setup 初始化日志配置
func Setup(cfg *Config) error {
	// 设置默认值
	if cfg.MaxSize == 0 {
		cfg.MaxSize = 100 // 默认100MB
	}
	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = 30 // 默认保留30个旧文件
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 7 // 默认保留7天
	}
	if cfg.Format == "" {
		cfg.Format = "json" // 默认使用JSON格式
	}
	if cfg.Level == "" {
		cfg.Level = "info" // 默认INFO级别
	}
	if cfg.Filename == "" {
		cfg.Filename = "./storage/logs/app.log" // 默认使用JSON格式
	}

	// 创建基本的encoder配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 设置日志级别
	var level zapcore.Level
	err := level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return fmt.Errorf("解析日志级别失败: %v", err)
	}

	// 确保日志目录存在
	if cfg.Filename != "" {
		logDir := filepath.Dir(cfg.Filename)
		if err := os.MkdirAll(logDir, 0o755); err != nil {
			return fmt.Errorf("创建日志目录失败: %v", err)
		}
	}

	// 创建Core
	var cores []zapcore.Core

	// 文件输出
	if cfg.Filename != "" {
		// 配置日志轮转
		hook := &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}

		var encoder zapcore.Encoder
		if cfg.Format == "json" {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		}

		cores = append(cores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(hook),
			level,
		))
	}

	// 控制台输出
	if cfg.Console {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		))
	}

	// 创建logger
	core := zapcore.NewTee(cores...)
	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// Field 定义日志字段
type Field = zapcore.Field

// String 创建字符串字段
func String(key, value string) Field {
	return zap.String(key, value)
}

// Int64 创建int64字段
func Int64(key string, value int64) Field {
	return zap.Int64(key, value)
}

// Err 创建error字段
func Err(err error) Field {
	return zap.Error(err)
}

// Info 记录info级别日志
func Info(msg string, fields ...Field) {
	log.Info(msg, fields...)
}

// Error 记录error级别日志
func Error(msg string, fields ...Field) {
	log.Error(msg, fields...)
}

// Debug 记录debug级别日志
func Debug(msg string, fields ...Field) {
	log.Debug(msg, fields...)
}

// Warn 记录warn级别日志
func Warn(msg string, fields ...Field) {
	log.Warn(msg, fields...)
}

// Fatal 记录fatal级别日志
func Fatal(msg string, fields ...Field) {
	log.Fatal(msg, fields...)
}

// Sync 同步日志缓冲
func Sync() error {
	return log.Sync()
}
