package tools

import (
	"context"
	"golang.org/x/exp/slog"
)

func NewMockLogger() *slog.Logger {
	return slog.New(NewMockHandler())
}

type MockHandler struct{}

func (m *MockHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return false
}

func (m *MockHandler) Handle(ctx context.Context, record slog.Record) error {
	return nil
}

func (m *MockHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return m
}

func (m *MockHandler) WithGroup(name string) slog.Handler {
	return m
}

func NewMockHandler() *MockHandler {
	return &MockHandler{}
}
