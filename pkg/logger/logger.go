package logger

import (
	"context"
)

//go:generate mockery --all --output ./mocks --outpkg loggermocks --with-expecter

type Logger interface {
	DebugLogger
	ErrorLogger
	InfoLogger
	WarnLogger

	With(args ...any) Logger
	WithGroup(name string) Logger
}

type DebugLogger interface {
	Debug(msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)
}

type ErrorLogger interface {
	Error(msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type InfoLogger interface {
	Info(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
}

type WarnLogger interface {
	Warn(msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
}
