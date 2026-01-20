package slogdiscard

import (
	"context"
	"log/slog"
)

// NewDiscardLogger returns a new logger that discards all log entries.
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

// Handle ignores the log record and returns nil.
func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

// WithAttrs returns the same handler as attributes are ignored.
func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup returns the same handler as groups are ignored.
func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}

// Enabled always returns false to skip log processing.
func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
