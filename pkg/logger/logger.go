package logger

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"iec104-slave/pkg/config"
	"io"
	"log/slog"
	"os"
)

// Init 初始化日志配置
func Init(config config.Logs) {
	level := slog.LevelInfo
	if config.Debug {
		level = slog.LevelDebug
	}

	file := "./logs/logs.txt"
	if len(config.File) > 0 {
		file = config.File
	}

	jack := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    100, // megabytes
		MaxBackups: 0,
		MaxAge:     1,     //days
		Compress:   false, // disabled by default
	}

	opts := &slog.HandlerOptions{
		AddSource: config.Debug,
		Level:     level,
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(io.MultiWriter(os.Stderr, jack), opts)))
}
