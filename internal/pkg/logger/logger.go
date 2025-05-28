package logger

import (
	"diprec_api/internal/config"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(cfg config.LoggingConfig) (*zap.Logger, error) {
	// Валидация конфигурации
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	// Настройка энкодера
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if cfg.JSONFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Настройка выходов
	var cores []zapcore.Core
	level := parseLogLevel(cfg.Level)

	// Добавляем вывод в файл если указан
	if cfg.LogFilePath != "" {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.LogFilePath,
			MaxSize:    cfg.RotationPolicy.MaxSize,
			MaxBackups: cfg.RotationPolicy.MaxBackups,
			MaxAge:     cfg.RotationPolicy.MaxAge,
			Compress:   true,
		})
		cores = append(cores, zapcore.NewCore(encoder, fileWriter, level))
	}

	// Всегда добавляем вывод в stdout
	stdoutWriter := zapcore.Lock(os.Stdout)
	cores = append(cores, zapcore.NewCore(encoder, stdoutWriter, level))

	// Создаем ядро
	core := zapcore.NewTee(cores...)

	// Настройка дополнительных опций
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return logger, nil
}

func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
