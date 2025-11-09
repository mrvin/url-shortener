package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type contextKey int

const (
	contextKeyRequestID contextKey = iota
	contextKeyUserName
)

const logFileMode = 0755

type Conf struct {
	FilePath string `yaml:"filepath"`
	Level    string `yaml:"level"`
}

type ContextHandler struct {
	slog.Handler
}

func Init(conf *Conf) (*os.File, error) {
	var level slog.Level

	if err := level.UnmarshalText([]byte(conf.Level)); err != nil {
		return nil, fmt.Errorf("getting level from text: %w", err)
	}

	var err error
	var logFile *os.File
	if conf.FilePath == "" {
		logFile = os.Stdout
	} else {
		logFile, err = os.OpenFile(conf.FilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, logFileMode)
		if err != nil {
			return nil, fmt.Errorf("failed open log file: %w", err)
		}
	}

	replaceAttr := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			t := a.Value.Any().(time.Time) //nolint:forcetypeassert
			a.Value = slog.StringValue(t.Format(time.StampNano))
		}
		return a
	}

	handler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		AddSource:   true,
		Level:       level,
		ReplaceAttr: replaceAttr,
	})

	logger := slog.New(ContextHandler{handler})
	slog.SetDefault(logger)

	return logFile, nil
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, contextKeyRequestID, requestID)
}

func WithUserName(ctx context.Context, userName string) context.Context {
	return context.WithValue(ctx, contextKeyUserName, userName)
}

func GetUserNameFromCtx(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", errors.New("ctx is nil")
	}
	userName, ok := ctx.Value(contextKeyUserName).(string)
	if !ok {
		return "", errors.New("no user name in ctx")
	}

	return userName, nil
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if requestID, ok := ctx.Value(contextKeyRequestID).(string); ok {
		r.Add("requestID", requestID)
	}
	if userName, ok := ctx.Value(contextKeyUserName).(string); ok {
		r.Add("userName", userName)
	}

	return h.Handler.Handle(ctx, r) //nolint:wrapcheck
}
