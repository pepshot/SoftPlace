package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/pepshot/SoftPlace/shared/config"
)

type Logger struct {
	*SlogWrapper
	closers []io.Closer
}

type SlogWrapper struct {
	*slog.Logger
}

func New(cfg config.LoggerConfig, serviceName string) (*Logger, error) {
	level := parseLevel(cfg.Level)

	writer, closers, err := buildWriter(cfg)
	if err != nil {
		return nil, err
	}

	options := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.SourceKey {
				source, ok := attr.Value.Any().(*slog.Source)
				if ok {
					attr.Value = slog.StringValue(formatSource(source))
				}
			}

			if attr.Key == slog.TimeKey {
				attr.Value = slog.StringValue(attr.Value.Time().Format("2006-01-02 15:04:05"))
			}

			return attr
		},
	}

	var handler slog.Handler

	switch strings.ToLower(cfg.Format) {
	case "json":
		handler = slog.NewJSONHandler(writer, options)
	case "text":
		handler = slog.NewTextHandler(writer, options)
	default:
		handler = slog.NewTextHandler(writer, options)
	}

	baseLogger := slog.New(handler).With(
		"service", serviceName,
	)

	return &Logger{
		SlogWrapper: &SlogWrapper{
			Logger: baseLogger,
		},
		closers: closers,
	}, nil
}

func (l *Logger) Close() error {
	var resultErr error

	for _, closer := range l.closers {
		if err := closer.Close(); err != nil {
			resultErr = err
		}
	}

	return resultErr
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func buildWriter(cfg config.LoggerConfig) (io.Writer, []io.Closer, error) {
	writers := make([]io.Writer, 0)
	closers := make([]io.Closer, 0)

	for _, output := range cfg.Outputs {
		switch strings.ToLower(strings.TrimSpace(output)) {
		case "stdout":
			writers = append(writers, os.Stdout)

		case "stderr":
			writers = append(writers, os.Stderr)

		case "file":
			file, err := openLogFile(cfg.FilePath)
			if err != nil {
				return nil, nil, err
			}

			writers = append(writers, file)
			closers = append(closers, file)

		default:
			return nil, nil, fmt.Errorf("unknown log output: %s", output)
		}
	}

	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	return io.MultiWriter(writers...), closers, nil
}

func openLogFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	return os.OpenFile(
		path,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
}

func formatSource(source *slog.Source) string {
	_, fileName, line, ok := runtime.Caller(0)
	_ = fileName
	_ = line
	_ = ok

	shortFile := source.File

	parts := strings.Split(source.File, string(os.PathSeparator))
	if len(parts) >= 2 {
		shortFile = filepath.Join(parts[len(parts)-2], parts[len(parts)-1])
	}

	return shortFile + ":" + strconv.Itoa(source.Line)
}
