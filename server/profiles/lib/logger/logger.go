package logger

import (
	"fmt"
	"golang.org/x/net/context"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
)

type HandlerOptions struct {
	TimestampFormat string
	Level           slog.Level
}

type CustomHandler struct {
	w    io.Writer
	opts HandlerOptions
	mu   *sync.Mutex
}

const (
	TimeFormat = "2006-01-02 15:04:05"
)

func NewCustomHandler(w io.Writer, opts *HandlerOptions) *CustomHandler {
	if opts == nil {
		opts = &HandlerOptions{
			TimestampFormat: TimeFormat,
			Level:           slog.LevelInfo,
		}
	}
	return &CustomHandler{
		w:    w,
		opts: *opts,
		mu:   &sync.Mutex{},
	}
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}
	h.mu.Lock()
	defer h.mu.Unlock()

	ts := r.Time.Format(h.opts.TimestampFormat)
	level := r.Level.String()
	msg := r.Message
	attrs := formatAttrs(r)

	if attrs == "" {
		fmt.Printf("%s | %s | %s\n", ts, level, msg)
		return nil
	}

	fmt.Printf("%s | %s | %s | %s\n", ts, level, msg, attrs)
	return nil
}

func formatAttrs(r slog.Record) string {
	var sb strings.Builder
	r.Attrs(func(attr slog.Attr) bool {
		sb.WriteString(fmt.Sprintf("%s=%v ", attr.Key, attr.Value))
		return true
	})
	return sb.String()
}

func (h *CustomHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level
}

func (h *CustomHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *CustomHandler) WithGroup(_ string) slog.Handler {
	return h
}

func SetupLogger(logLevel string) *slog.Logger {
	var log *slog.Logger
	logLevel = strings.ToLower(logLevel)

	if logLevel == "dev" {
		log = slog.New(
			NewCustomHandler(os.Stdout, &HandlerOptions{TimestampFormat: TimeFormat, Level: slog.LevelDebug}),
		)
		return log
	}

	log = slog.New(
		NewCustomHandler(os.Stdout, &HandlerOptions{TimestampFormat: TimeFormat, Level: slog.LevelWarn}),
	)

	return log
}
