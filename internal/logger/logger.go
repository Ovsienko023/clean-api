package logger

import (
	"log/slog"
	"os"
)

const (
	FormatJson string = "json"
	FormatText string = "text"

	OutputPathStdOut string = "stdout"
)

type Config struct {
	Level      slog.Level // Уровень логирования
	OutputPath string     // Путь к файлу для логирования (если пусто – вывод в stdout)
	Format     string     // Формат логирования: "json" или "text"
}

// New создаёт и возвращает настроенный логгер.
func New(cfg Config) (*slog.Logger, error) {
	var output *os.File
	if cfg.OutputPath == OutputPathStdOut {
		output = os.Stdout
	} else {
		var err error
		output, err = os.OpenFile(cfg.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: cfg.Level,
	}

	switch cfg.Format {
	case FormatJson:
		handler = slog.NewJSONHandler(output, opts)
	case FormatText:
		handler = slog.NewTextHandler(output, opts)
	default:
		handler = slog.NewTextHandler(output, opts)
	}

	logger := slog.New(handler)
	return logger, nil
}
