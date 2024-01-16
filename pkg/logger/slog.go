package logger

import "log/slog"

type SlogWrapper struct {
	*slog.Logger
}

func (wr *SlogWrapper) With(args ...any) Logger {
	return &SlogWrapper{Logger: wr.Logger.With(args...)}
}

func (wr *SlogWrapper) WithGroup(name string) Logger {
	return &SlogWrapper{Logger: wr.Logger.WithGroup(name)}
}
