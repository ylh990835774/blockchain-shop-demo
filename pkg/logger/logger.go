package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Level      string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

type Logger struct {
	*zap.Logger
	config *Config
}

func NewLogger(config *Config) (*Logger, error) {
	// 创建基础目录
	logDir := filepath.Dir(config.Filename)
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, fmt.Errorf("create log directory failed: %v", err)
	}

	// 生成当前日期的日志文件名
	now := time.Now()
	baseFileName := filepath.Base(config.Filename)
	ext := filepath.Ext(baseFileName)
	name := baseFileName[:len(baseFileName)-len(ext)]
	dailyFileName := fmt.Sprintf("%s.%s%s", name, now.Format("2006-01-02"), ext)
	dailyFilePath := filepath.Join(logDir, dailyFileName)

	// 配置 lumberjack
	writer := &lumberjack.Logger{
		Filename:   dailyFilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	// 创建软链接指向当前日志文件
	if err := createSymlink(dailyFileName, config.Filename); err != nil {
		return nil, fmt.Errorf("create symlink failed: %v", err)
	}

	// 配置 zap
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
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

	atomicLevel := zap.NewAtomicLevel()
	if err := atomicLevel.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("parse log level failed: %v", err)
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		atomicLevel,
	)

	logger := &Logger{
		Logger: zap.New(core, zap.AddCaller()),
		config: config,
	}

	// 启动定时器，每天凌晨更新日志文件
	go logger.rotateDaily()

	return logger, nil
}

func createSymlink(dailyFileName, symlinkPath string) error {
	// 删除已存在的软链接
	_ = os.Remove(symlinkPath)

	// 获取工作目录
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory failed: %v", err)
	}

	// 计算相对路径
	symlinkDir := filepath.Dir(symlinkPath)
	if err := os.MkdirAll(symlinkDir, 0o755); err != nil {
		return fmt.Errorf("create symlink directory failed: %v", err)
	}

	// 创建软链接（使用相对路径）
	if err := os.Chdir(symlinkDir); err != nil {
		return fmt.Errorf("change directory failed: %v", err)
	}
	defer os.Chdir(wd)

	if err := os.Symlink(dailyFileName, filepath.Base(symlinkPath)); err != nil {
		return fmt.Errorf("create symlink failed: %v", err)
	}

	return nil
}

func (l *Logger) rotateDaily() {
	for {
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		duration := next.Sub(now)

		timer := time.NewTimer(duration)
		<-timer.C

		// 生成新的日志文件名
		baseFileName := filepath.Base(l.config.Filename)
		ext := filepath.Ext(baseFileName)
		name := baseFileName[:len(baseFileName)-len(ext)]
		dailyFileName := fmt.Sprintf("%s.%s%s", name, time.Now().Format("2006-01-02"), ext)

		// 更新软链接
		_ = createSymlink(dailyFileName, l.config.Filename)
	}
}
