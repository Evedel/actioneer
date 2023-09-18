package logging

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"
)

// Init initializes logger with log level and writer
// log_level: debug | info | warn | error
// writer: if nil, will use os.Stdout
func Init(log_level string, writer io.Writer) (e error) {
	e = nil

	leveler := new(slog.LevelVar)
	if writer == nil {
		writer = os.Stdout
	}
	h := slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: leveler})
	slog.SetDefault(slog.New(h))

	switch log_level {
	case "debug":
		leveler.Set(slog.LevelDebug)
	case "info":
		leveler.Set(slog.LevelInfo)
	case "warn":
		leveler.Set(slog.LevelWarn)
	case "error":
		leveler.Set(slog.LevelError)
	default:
		slog.Error("wrong value in --log-level=" + log_level)
		e = errors.New("wrong value in --log-level=" + log_level)
	}
	slog.Info("--log-level=" + strings.ToLower(leveler.Level().String()))

	return
}
